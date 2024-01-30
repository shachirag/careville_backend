package physiotherapist

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

// @Summary Get all services
// @Description Get all investigations
// @Tags physiotherapist
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} services.PhysiotherapistServicesResDto
// @Router /provider/services/get-physiotherapist-services [get]
func GetPhysiotherapistServices(c *fiber.Ctx) error {
	ctx := context.Background()

	var service entity.ServiceEntity

	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"_id": providerData.ProviderId,
	}

	projection := bson.M{
		"physiotherapist.serviceAndSchedule.id":          1,
		"physiotherapist.serviceAndSchedule.name":        1,
		"physiotherapist.serviceAndSchedule.serviceFees": 1,
		"physiotherapist.serviceAndSchedule.slots": bson.M{
			"startTime": 1,
			"endTime":   1,
			"days":      1,
		},
	}

	findOptions := options.FindOne().SetProjection(projection)

	err := serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.PhysiotherapistServicesResDto{
				Status:  false,
				Message: "service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.PhysiotherapistServicesResDto{
			Status:  false,
			Message: "Failed to fetch services from MongoDB: " + err.Error(),
		})
	}

	serviceData := make([]services.PhysiotherapistServiceRes, 0)
	if service.Physiotherapist != nil && len(service.Physiotherapist.ServiceAndSchedule) > 0 {
		for _, service := range service.Physiotherapist.ServiceAndSchedule {
			// Create a new slice to hold slots data for this particular service
			var slotData []services.Slots
			for _, slots := range service.Slots {
				slotData = append(slotData, services.Slots{
					StartTime: slots.StartTime,
					EndTime:   slots.EndTime,
					Days:      slots.Days,
				})
			}
			serviceData = append(serviceData, services.PhysiotherapistServiceRes{
				Id:          service.Id,
				ServiceFees: service.ServiceFees,
				Name:        service.Name,
				Slots:       slotData,
			})
		}

	}

	if len(serviceData) == 0 {
		return c.Status(fiber.StatusOK).JSON(services.PhysiotherapistServicesResDto{
			Status:  false,
			Message: "No service data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(services.PhysiotherapistServicesResDto{
		Status:  true,
		Message: "services retrieved successfully",
		Data:    serviceData,
	})
}
