package fitnessCenter

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

// @Summary Update other service info
// @Description Update other service info
// @Tags fitnessCenter
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param provider body services.UpdateOtherServiceReqDto true "Update data of trainer"
// @Produce json
// @Success 200 {object} services.UpdateTrainerResDto
// @Router /provider/services/update-fitnessCenter-other-service-info/{otherServiceId} [put]
func UpdateOtherServiceInfo(c *fiber.Ctx) error {

	var (
		serviceColl = database.GetCollection("service")
		data        services.UpdateOtherServiceReqDto
		provider    entity.ServiceEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(services.UpdateTrainerResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	otherServiceId := c.Params("otherServiceId")
	otherServiceObjID, err := primitive.ObjectIDFromHex(otherServiceId)

	if err != nil {
		return c.Status(400).JSON(services.UpdateDoctorImageResDto{
			Status:  false,
			Message: "invalid objectId " + err.Error(),
		})
	}

	filter := bson.M{
		"_id": providerData.ProviderId,
		"fitnessCenter.additionalServices": bson.M{
			"$elemMatch": bson.M{
				"id": otherServiceObjID,
			},
		},
	}

	err = serviceColl.FindOne(ctx, filter).Decode(&provider)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.UpdateTrainerResDto{
				Status:  false,
				Message: "other service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateTrainerResDto{
			Status:  false,
			Message: "Failed to fetch other service from MongoDB: " + err.Error(),
		})
	}

	update := bson.M{
		"$set": bson.M{
			"fitnessCenter.additionalServices.$.name":        data.Name,
			"fitnessCenter.additionalServices.$.information": data.Information,
			"updatedAt": time.Now().UTC(),
		},
	}

	// Execute the update operation
	updateRes, err := serviceColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateTrainerResDto{
			Status:  false,
			Message: "Failed to update other service data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(services.UpdateTrainerResDto{
			Status:  false,
			Message: "other service not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(services.UpdateTrainerResDto{
		Status:  true,
		Message: "other service data updated successfully",
	})
}
