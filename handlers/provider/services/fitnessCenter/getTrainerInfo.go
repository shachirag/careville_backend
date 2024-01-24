package fitnessCenter

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get trainer info
// @Description Get trainer info
// @Tags fitnessCenter
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param trainerId path string true "trainer ID"
// @Produce json
// @Success 200 {object} services.GetTrainerResDto
// @Router /provider/services/get-trainer-info/{trainerId} [get]
func GetTrainerInfo(c *fiber.Ctx) error {

	var service entity.ServiceEntity

	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	trainerId := c.Params("trainerId")
	trainerObjID, err := primitive.ObjectIDFromHex(trainerId)

	if err != nil {
		return c.Status(400).JSON(services.GetTrainerResDto{
			Status:  false,
			Message: "invalid objectId " + err.Error(),
		})
	}

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"_id": providerData.ProviderId,
		"fitnessCenter.trainers": bson.M{
			"$elemMatch": bson.M{
				"id": trainerObjID,
			},
		},
	}

	projection := bson.M{
		"fitnessCenter.trainers.id":          1,
		"fitnessCenter.trainers.name":        1,
		"fitnessCenter.trainers.price":       1,
		"fitnessCenter.trainers.information": 1,
		"fitnessCenter.trainers.category":    1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.GetTrainerResDto{
				Status:  false,
				Message: "service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.GetTrainerResDto{
			Status:  false,
			Message: "Failed to fetch service from MongoDB: " + err.Error(),
		})
	}

	if service.FitnessCenter == nil {
		return c.Status(fiber.StatusNotFound).JSON(services.GetTrainerResDto{
			Status:  false,
			Message: "No trainer information found for the service",
		})
	}

	var trainersRes services.TrainerRes

	for _, trainers := range service.FitnessCenter.Trainers {
		if trainers.Id == trainerObjID {
			trainerRes := services.TrainerRes{
				Id:          trainers.Id,
				Name:        trainers.Name,
				Category:    trainers.Category,
				Information: trainers.Information,
				Price:       trainers.Price,
			}

			trainersRes = trainerRes
			break
		}
	}

	return c.Status(fiber.StatusOK).JSON(services.GetTrainerResDto{
		Status:  true,
		Message: "trainer retrieved successfully",
		Data:    trainersRes,
	})
}
