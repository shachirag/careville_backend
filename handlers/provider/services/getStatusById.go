package services

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Fetch status By ID
// @Description Fetch status By ID
// @Tags services
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} services.StatusRes
// @Router /provider/services/get-status/{id} [get]
func FetchStatusById(c *fiber.Ctx) error {

	var service entity.ServiceEntity

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	serviceColl := database.GetCollection("service")

	err := serviceColl.FindOne(ctx, bson.M{"_id": providerData.ProviderId}).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.StatusRes{
				Status:  false,
				Message: "status not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.StatusRes{
			Status:  false,
			Message: "Failed to fetch status from MongoDB: " + err.Error(),
		})
	}

	statusRes := services.StatusRespDto{
		Id:            service.Id,
		ServiceStatus: service.ServiceStatus,
	}

	return c.Status(fiber.StatusOK).JSON(services.StatusRes{
		Status:  true,
		Message: "status retrieved successfully",
		Data:    statusRes,
	})
}
