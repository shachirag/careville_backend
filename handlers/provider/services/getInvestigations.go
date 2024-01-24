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

// @Summary GetAllDoctors
// @Description GetAllDoctors
// @Tags services
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} services.InvestigationResDto
// @Router /provider/services/get-investigations [get]
func GetInvestigations(c *fiber.Ctx) error {
	ctx := context.TODO()

	var service entity.ServiceEntity

	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"_id": providerData.ProviderId,
	}

	projection := bson.M{
		"laboratory.investigations.id":          1,
		"laboratory.investigations.name":        1,
		"laboratory.investigations.type":        1,
		"laboratory.investigations.information": 1,
		"laboratory.investigations.price":       1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	err := serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.InvestigationResDto{
				Status:  false,
				Message: "investigation not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.InvestigationResDto{
			Status:  false,
			Message: "Failed to fetch investigation from MongoDB: " + err.Error(),
		})
	}

	investigationData := make([]services.InvestigationRes, 0)
	if service.Laboratory != nil && len(service.Laboratory.Investigations) > 0 {
		for _, investigation := range service.Laboratory.Investigations {
			investigationData = append(investigationData, services.InvestigationRes{
				Id:          investigation.Id,
				Type:        investigation.Type,
				Name:        investigation.Name,
				Information: investigation.Information,
				Price:       investigation.Price,
			})
		}
	}
	if len(investigationData) == 0 {
		return c.Status(fiber.StatusOK).JSON(services.InvestigationResDto{
			Status:  false,
			Message: "No investigation data found.",
		})
	}
	return c.Status(fiber.StatusOK).JSON(services.InvestigationResDto{
		Status:  true,
		Message: "investigations retrieved successfully",
		Data:    investigationData,
	})
}
