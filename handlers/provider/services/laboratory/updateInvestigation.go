package laboratory

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Update investigation info
// @Description Update investigation info
// @Tags laboratory
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param investigationId path string true "investigation ID"
// @Param provider body services.UpdateInvestigationReqDto true "Update data of investigation"
// @Produce json
// @Success 200 {object} services.UpdateInvestigationResDto
// @Router /provider/services/update-investigation-info/{investigationId} [put]
func UpdateinvestigationInfo(c *fiber.Ctx) error {

	var (
		serviceColl = database.GetCollection("service")
		data        services.UpdateInvestigationReqDto
		provider    entity.ServiceEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(services.UpdateInvestigationResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	investigationId := c.Params("investigationId")
	investigationObjID, err := primitive.ObjectIDFromHex(investigationId)

	if err != nil {
		return c.Status(400).JSON(services.UpdateDoctorImageResDto{
			Status:  false,
			Message: "invalid objectId " + err.Error(),
		})
	}

	filter := bson.M{
		"_id": providerData.ProviderId,
		"laboratory.investigations": bson.M{
			"$elemMatch": bson.M{
				"id": investigationObjID,
			},
		},
	}

	err = serviceColl.FindOne(ctx, filter).Decode(&provider)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.UpdateInvestigationResDto{
				Status:  false,
				Message: "Investigation not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateInvestigationResDto{
			Status:  false,
			Message: "Failed to fetch investigation from MongoDB: " + err.Error(),
		})
	}

	update := bson.M{
		"$set": bson.M{
			"laboratory.investigations.$.type":        data.Type,
			"laboratory.investigations.$.name":        data.Name,
			"laboratory.investigations.$.information": data.Information,
			"laboratory.investigations.$.price":       data.Price,
			"updatedAt":                               time.Now().UTC(),
		},
	}

	// Execute the update operation
	updateRes, err := serviceColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateInvestigationResDto{
			Status:  false,
			Message: "Failed to update investigation data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(services.UpdateInvestigationResDto{
			Status:  false,
			Message: "Investigation not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(services.UpdateInvestigationResDto{
		Status:  true,
		Message: "Investigation data updated successfully",
	})
}
