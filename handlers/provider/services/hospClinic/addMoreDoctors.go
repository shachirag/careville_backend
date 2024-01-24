package hospClinic

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Add AddMoreDoctors
// @Tags hospClinic
// @Description Add AddMoreDoctors
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param  provider body services.MoreDoctorReqDto true "add AddMoreDoctors"
// @Produce json
// @Success 200 {object} services.UpdateDoctorResDto
// @Router /provider/services/add-more-doctor [post]
func AddMoreDoctors(c *fiber.Ctx) error {
	var (
		servicesColl = database.GetCollection("service")
		data         services.MoreDoctorReqDto
		provider     entity.ServiceEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(services.UpdateDoctorResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	filter := bson.M{
		"_id": providerData.ProviderId,
	}

	err = servicesColl.FindOne(ctx, filter).Decode(&provider)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.UpdateDoctorImageResDto{
				Status:  false,
				Message: "provider not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateDoctorImageResDto{
			Status:  false,
			Message: "Failed to fetch provider from MongoDB: " + err.Error(),
		})
	}

	var schedule []entity.Schedule
	for _, inv := range data.Schedule {
		convertedInv := entity.Schedule{
			StartTime: inv.StartTime,
			EndTime:   inv.EndTime,
			Days:      inv.Days,
		}
		schedule = append(schedule, convertedInv)
	}

	moreDoctor := []entity.Doctor{
		{
			Id:         primitive.NewObjectID(),
			Name:       data.Name,
			Speciality: data.Speciality,
			Schedule:   schedule,
		},
	}

	update := bson.M{
		"$push": bson.M{"hospClinic.doctor": bson.M{"$each": moreDoctor}},
	}

	updateRes, err := servicesColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.HospitalClinicResDto{
			Status:  false,
			Message: "Failed to update provider data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(services.HospitalClinicResDto{
			Status:  false,
			Message: "provider not found",
		})
	}

	hospClinicRes := services.HospitalClinicResDto{
		Status:  true,
		Message: "doctor added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(hospClinicRes)
}
