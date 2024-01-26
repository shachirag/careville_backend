package fitnessCenter

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get other services
// @Description Get other services
// @Tags fitnessCenter
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} services.GetFitnessOtherServicesResDto
// @Router /provider/services/get-fitness-other-services [get]
func GetOtherServices(c *fiber.Ctx) error {

	var service entity.ServiceEntity

	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"_id": providerData.ProviderId,
	}

	projection := bson.M{
		"fitnessCenter.additionalServices.id":          1,
		"fitnessCenter.additionalServices.name":        1,
		"fitnessCenter.additionalServices.information": 1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	err := serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.GetFitnessOtherServicesResDto{
				Status:  false,
				Message: "other service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.GetFitnessOtherServicesResDto{
			Status:  false,
			Message: "Failed to fetch other service from MongoDB: " + err.Error(),
		})
	}

	serviceData := make([]services.FitnessOtherServiceRes, 0)
	if service.FitnessCenter != nil && len(service.FitnessCenter.AdditionalServices) > 0 {
		for _, service := range service.FitnessCenter.AdditionalServices {
			serviceData = append(serviceData, services.FitnessOtherServiceRes{
				Id:          service.Id,
				Name:        service.Name,
				Information: service.Information,
			})
		}
	}

	if len(serviceData) == 0 {
		return c.Status(fiber.StatusOK).JSON(services.GetFitnessOtherServicesResDto{
			Status:  false,
			Message: "No other service data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(services.GetFitnessOtherServicesResDto{
		Status:  true,
		Message: "other service retrieved successfully",
		Data:    serviceData,
	})
}
