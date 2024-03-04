package doctorProfession

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	services "careville_backend/dto/provider/services"
	"careville_backend/entity"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Update Professional details
// @Description Update Professional details
// @Tags doctorProfession
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param provider body services.UpdateDoctorProfessionProfessionalInfoReqDto true "Update data of provider"
// @Produce json
// @Success 200 {object} services.UpdateDoctorProfessionProfessionalInfoResDto
// @Router /provider/services/edit-doctorProfession-professional-info [put]
func UpdateDoctorProfessionDetails(c *fiber.Ctx) error {

	var (
		serviceColl = database.GetCollection("service")
		data        services.UpdateDoctorProfessionProfessionalInfoReqDto
		provider    entity.ServiceEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(services.UpdateDoctorProfessionProfessionalInfoResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	filter := bson.M{"_id": providerData.ProviderId}
	err = serviceColl.FindOne(ctx, filter).Decode(&provider)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.UpdateDoctorProfessionProfessionalInfoResDto{
				Status:  false,
				Message: "provider not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateDoctorProfessionProfessionalInfoResDto{
			Status:  false,
			Message: "Failed to fetch provider from MongoDB: " + err.Error(),
		})
	}

	update := bson.M{}

	if provider.Doctor != nil {
		update = bson.M{"$set": bson.M{
			"doctor.additionalServices.qualifications": data.Qualifications,
			"doctor.additionalServices.speciality":     data.Speciality,
			"doctor.schedule.consultationFees":         data.ConsultingFees,
			"updatedAt":                                time.Now().UTC(),
		},
		}
	}

	opts := options.Update().SetUpsert(true)

	session, err := database.GetMongoClient().StartSession()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateDoctorProfessionProfessionalInfoResDto{
			Status:  false,
			Message: "Failed to start session",
		})
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		updateRes, err := serviceColl.UpdateOne(sessCtx, filter, update, opts)
		if err != nil {
			return nil, err
		}

		if updateRes.MatchedCount == 0 {
			return nil, mongo.ErrNoDocuments
		}

		appointmentUpdate := bson.M{"$set": bson.M{
			"doctor.information.speciality": data.Speciality,
		}}

		_, err = database.GetCollection("appointment").UpdateMany(sessCtx, bson.M{"serviceId": providerData.ProviderId}, appointmentUpdate)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateDoctorProfessionProfessionalInfoResDto{
			Status:  false,
			Message: "Failed to update appointment data: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(services.UpdateDoctorProfessionProfessionalInfoResDto{
		Status:  true,
		Message: "Professional Details data updated successfully",
	})
}
