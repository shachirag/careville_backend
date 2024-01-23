package services

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary get other services for provider
// @Description get other services for provider
// @Tags services
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} services.OtherServicesResDto
// @Router /provider/services/get-other-services [get]
func GetOtherServices(c *fiber.Ctx) error {
	ctx := context.TODO()

	var service entity.ServiceEntity

	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"_id": providerData.ProviderId,
	}

	projection := bson.M{
		"hospClinic.otherServices": 1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	err := serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.OtherServicesResDto{
				Status:  false,
				Message: "other service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.OtherServicesResDto{
			Status:  false,
			Message: "Failed to fetch other service from MongoDB: " + err.Error(),
		})
	}

	if service.HospClinic == nil {
		return c.Status(fiber.StatusOK).JSON(services.OtherServicesResDto{
			Status:  false,
			Message: "HospClinic information not found.",
		})
	}

	response := services.OtherServicesResDto{
		Status:  true,
		Message: "Other services retrieved successfully",
		Data: services.OtherServicesRes{
			OtherServices: service.HospClinic.OtherServices,
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
