package medicalLabScientist

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
// @Tags medicalLabScientist
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} services.MedicalLabScientistServicesResDto
// @Router /provider/services/get-medicalLabScientist-services [get]
func GetMedicalLabScientistServices(c *fiber.Ctx) error {
	ctx := context.Background()

	var service entity.ServiceEntity

	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"_id": providerData.ProviderId,
	}

	projection := bson.M{
		"medicalLabScientist.serviceAndSchedule.id":          1,
		"medicalLabScientist.serviceAndSchedule.name":        1,
		"medicalLabScientist.serviceAndSchedule.serviceFees": 1,
		"medicalLabScientist.serviceAndSchedule.slots": bson.M{
			"startTime": 1,
			"endTime":   1,
			"days":      1,
		},
	}

	findOptions := options.FindOne().SetProjection(projection)

	err := serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.MedicalLabScientistServicesResDto{
				Status:  false,
				Message: "service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.MedicalLabScientistServicesResDto{
			Status:  false,
			Message: "Failed to fetch services from MongoDB: " + err.Error(),
		})
	}

	serviceData := make([]services.MedicalLabScientistServiceRes, 0)
	if service.MedicalLabScientist != nil && len(service.MedicalLabScientist.ServiceAndSchedule) > 0 {
		for _, service := range service.MedicalLabScientist.ServiceAndSchedule {
			// Create a new slice to hold slots data for this particular service
			var slotData []services.Slots
			for _, slots := range service.Slots {
				slotData = append(slotData, services.Slots{
					StartTime: slots.StartTime,
					EndTime:   slots.EndTime,
					Days:      slots.Days,
				})
			}
			serviceData = append(serviceData, services.MedicalLabScientistServiceRes{
				Id:          service.Id,
				ServiceFees: service.ServiceFees,
				Name:        service.Name,
				Slots:       slotData,
			})
		}

	}

	if len(serviceData) == 0 {
		return c.Status(fiber.StatusOK).JSON(services.MedicalLabScientistServicesResDto{
			Status:  false,
			Message: "No service data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(services.MedicalLabScientistServicesResDto{
		Status:  true,
		Message: "services retrieved successfully",
		Data:    serviceData,
	})
}
