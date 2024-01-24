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

// @Summary Update trainer info
// @Description Update trainer info
// @Tags fitnessCenter
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param trainerId path string true "trainer ID"
// @Param provider body services.UpdateTrainerReqDto true "Update data of trainer"
// @Produce json
// @Success 200 {object} services.UpdateTrainerResDto
// @Router /provider/services/update-trainer-info/{trainerId} [put]
func UpdateTrainerInfo(c *fiber.Ctx) error {

	var (
		serviceColl = database.GetCollection("service")
		data        services.UpdateTrainerReqDto
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

	trainerId := c.Params("trainerId")
	trainerObjID, err := primitive.ObjectIDFromHex(trainerId)

	if err != nil {
		return c.Status(400).JSON(services.UpdateDoctorImageResDto{
			Status:  false,
			Message: "invalid objectId " + err.Error(),
		})
	}

	filter := bson.M{
		"_id": providerData.ProviderId,
		"fitnessCenter.trainers": bson.M{
			"$elemMatch": bson.M{
				"id": trainerObjID,
			},
		},
	}

	err = serviceColl.FindOne(ctx, filter).Decode(&provider)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.UpdateTrainerResDto{
				Status:  false,
				Message: "trainer not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateTrainerResDto{
			Status:  false,
			Message: "Failed to fetch trainer from MongoDB: " + err.Error(),
		})
	}

	update := bson.M{
		"$set": bson.M{
			"fitnessCenter.trainers.$.category":    data.Category,
			"fitnessCenter.trainers.$.name":        data.Name,
			"fitnessCenter.trainers.$.information": data.Information,
			"fitnessCenter.trainers.$.price":       data.Price,
			"updatedAt":                            time.Now().UTC(),
		},
	}

	// Execute the update operation
	updateRes, err := serviceColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateTrainerResDto{
			Status:  false,
			Message: "Failed to update trainer data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(services.UpdateTrainerResDto{
			Status:  false,
			Message: "trainer not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(services.UpdateTrainerResDto{
		Status:  true,
		Message: "trainer data updated successfully",
	})
}
