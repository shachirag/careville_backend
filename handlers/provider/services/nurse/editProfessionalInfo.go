package nurse

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
// @Tags nurse
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param provider body services.UpdateNurseProfessionalInfoReqDto true "Update data of provider"
// @Produce json
// @Success 200 {object} services.UpdateNurseProfessionalInfoResDto
// @Router /provider/services/edit-nurse-professional-info [put]
func UpdateNurseDetails(c *fiber.Ctx) error {

	var (
		serviceColl = database.GetCollection("service")
		data        services.UpdateNurseProfessionalInfoReqDto
		provider    entity.ServiceEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(services.UpdateNurseProfessionalInfoResDto{
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
			return c.Status(fiber.StatusNotFound).JSON(services.UpdateNurseProfessionalInfoResDto{
				Status:  false,
				Message: "provider not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateNurseProfessionalInfoResDto{
			Status:  false,
			Message: "Failed to fetch provider from MongoDB: " + err.Error(),
		})
	}

	update := bson.M{}

	if provider.Nurse != nil {
		update = bson.M{"$set": bson.M{
			"nurse.professionalDetails.qualifications": data.Qualifications,
			"updatedAt": time.Now().UTC(),
		},
		}
	}

	opts := options.Update().SetUpsert(true)
	// Execute the update operation
	updateRes, err := serviceColl.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateNurseProfessionalInfoResDto{
			Status:  false,
			Message: "Failed to update provider data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(services.UpdateNurseProfessionalInfoResDto{
			Status:  false,
			Message: "provider not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(services.UpdateNurseProfessionalInfoResDto{
		Status:  true,
		Message: "provider data updated successfully",
	})
}
