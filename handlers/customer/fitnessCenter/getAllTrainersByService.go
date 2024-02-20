package fitnessCenter

import (
	"careville_backend/database"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get all trainers
// @Description Get all trainers
// @Tags customer fitnessCenter
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param serviceId query string true "service ID"
// @Produce json
// @Success 200 {object} services.TrainerResDto
// @Router /customer/healthFacility/get-all-trainers [get]
func GetAllTrainers(c *fiber.Ctx) error {

	var service entity.ServiceEntity

	serviceId := c.Query("serviceId")

	if serviceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(services.TrainerResDto{
			Status:  false,
			Message: "service Id is mandatory",
		})
	}

	serviceObjectID, err := primitive.ObjectIDFromHex(serviceId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.TrainerResDto{
			Status:  false,
			Message: "Invalid ID format",
		})
	}

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"_id": serviceObjectID,
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
			return c.Status(fiber.StatusNotFound).JSON(services.TrainerResDto{
				Status:  false,
				Message: "service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.TrainerResDto{
			Status:  false,
			Message: "Failed to fetch service from MongoDB: " + err.Error(),
		})
	}

	if service.FitnessCenter == nil {
		return c.Status(fiber.StatusNotFound).JSON(services.TrainerResDto{
			Status:  false,
			Message: "No FitnessCenter information found for the service",
		})
	}

	trainersByCategory := make(map[string][]services.TrainerRes)

	if service.FitnessCenter != nil && len(service.FitnessCenter.Trainers) > 0 {
		for _, trainer := range service.FitnessCenter.Trainers {
			trainerRes := services.TrainerRes{
				Id:          trainer.Id,
				Name:        trainer.Name,
				Category:    trainer.Category,
				Information: trainer.Information,
				Price:       trainer.Price,
			}
			trainersByCategory[trainer.Category] = append(trainersByCategory[trainer.Category], trainerRes)
		}
	}

	var response []services.SpecialityTrainerRes

	for category, trainers := range trainersByCategory {
		response = append(response, services.SpecialityTrainerRes{
			Category: category,
			Trainers: trainers,
		})
	}

	return c.Status(fiber.StatusOK).JSON(services.TrainerResDto{
		Status:  true,
		Message: "trainers retrieved successfully",
		Data:    response,
	})
}
