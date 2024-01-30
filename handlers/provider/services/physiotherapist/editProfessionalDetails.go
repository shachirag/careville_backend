package physiotherapist

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	providerAuth "careville_backend/dto/provider/providerAuth"
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
// @Tags physiotherapist
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param provider body services.UpdatePhysiotherapistProfessionalInfoReqDto true "Update data of provider"
// @Produce json
// @Success 200 {object} services.UpdatePhysiotherapistProfessionalInfoResDto
// @Router /provider/services/edit-physiotherapist-professional-info [put]
func UpdatePhysiotherapistDetails(c *fiber.Ctx) error {

	var (
		serviceColl = database.GetCollection("service")
		data        services.UpdatePhysiotherapistProfessionalInfoReqDto
		provider    entity.ServiceEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(services.UpdatePhysiotherapistProfessionalInfoResDto{
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
			return c.Status(fiber.StatusNotFound).JSON(providerAuth.GetProviderResDto{
				Status:  false,
				Message: "provider not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.GetProviderResDto{
			Status:  false,
			Message: "Failed to fetch provider from MongoDB: " + err.Error(),
		})
	}

	update := bson.M{}

	if provider.Physiotherapist != nil {
		update = bson.M{"$set": bson.M{
			"physiotherapist.professionalDetails.qualifications": data.Qualifications,
			"updatedAt": time.Now().UTC(),
		},
		}
	}

	opts := options.Update().SetUpsert(true)
	// Execute the update operation
	updateRes, err := serviceColl.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdatePhysiotherapistProfessionalInfoResDto{
			Status:  false,
			Message: "Failed to update provider data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(services.UpdatePhysiotherapistProfessionalInfoResDto{
			Status:  false,
			Message: "provider not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(services.UpdatePhysiotherapistProfessionalInfoResDto{
		Status:  true,
		Message: "provider data updated successfully",
	})
}
