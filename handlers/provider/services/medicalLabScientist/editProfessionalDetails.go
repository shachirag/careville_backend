package medicalLabScientist

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Update provider
// @Description Update provider
// @Tags medicalLabScientist
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param provider body services.UpdateMedicalLabScientistProfessionalInfoReqDto true "Update data of provider"
// @Produce json
// @Success 200 {object} services.UpdateMedicalLabScientistProfessionalInfoResDto
// @Router /provider/services/edit-medicalLabScientist-professional-info [put]
func UpdateMedicalLabScientistDetails(c *fiber.Ctx) error {

	var (
		serviceColl = database.GetCollection("service")
		data        services.UpdateMedicalLabScientistProfessionalInfoReqDto
		provider    entity.ServiceEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(services.UpdateMedicalLabScientistProfessionalInfoResDto{
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
			return c.Status(fiber.StatusNotFound).JSON(services.UpdateMedicalLabScientistProfessionalInfoResDto{
				Status:  false,
				Message: "provider not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateMedicalLabScientistProfessionalInfoResDto{
			Status:  false,
			Message: "Failed to fetch provider from MongoDB: " + err.Error(),
		})
	}

	update := bson.M{}

	if provider.MedicalLabScientist != nil {
		update = bson.M{"$set": bson.M{
			"medicalLabScientist.professionalDetails.qualification": data.Qualifications,
			"medicalLabScientist.professionalDetails.department":    data.Department,
			"updatedAt": time.Now().UTC(),
		},
		}
	}

	opts := options.Update().SetUpsert(true)
	// Execute the update operation
	updateRes, err := serviceColl.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateMedicalLabScientistProfessionalInfoResDto{
			Status:  false,
			Message: "Failed to update provider data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(services.UpdateMedicalLabScientistProfessionalInfoResDto{
			Status:  false,
			Message: "provider not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(services.UpdateMedicalLabScientistProfessionalInfoResDto{
		Status:  true,
		Message: "provider data updated successfully",
	})
}
