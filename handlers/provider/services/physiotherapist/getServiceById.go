package physiotherapist

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get service info
// @Description Get service info
// @Tags physiotherapist
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param serviceId path string true "serviceId"
// @Produce json
// @Success 200 {object} services.DoctorResDto
// @Router /provider/services/get-physiotherapist-service-info/{serviceId} [get]
func GetServiceInfo(c *fiber.Ctx) error {
	ctx := context.TODO()

	var service entity.ServiceEntity

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	serviceId := c.Params("serviceId")
	serviceObjID, err := primitive.ObjectIDFromHex(serviceId)

	if err != nil {
		return c.Status(400).JSON(services.GetPhysiotherapistServicesResDto{
			Status:  false,
			Message: "invalid objectId " + err.Error(),
		})
	}

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"_id": providerData.ProviderId,
		"physiotherapist.serviceAndSchedule": bson.M{
			"$elemMatch": bson.M{
				"id": serviceObjID,
			},
		},
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

	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.GetPhysiotherapistServicesResDto{
				Status:  false,
				Message: "service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.GetPhysiotherapistServicesResDto{
			Status:  false,
			Message: "Failed to fetch service from MongoDB: " + err.Error(),
		})
	}

	if service.Physiotherapist == nil {
		return c.Status(fiber.StatusNotFound).JSON(services.GetPhysiotherapistServicesResDto{
			Status:  false,
			Message: "No service information found for the service",
		})
	}

	var servicesRes services.PhysiotherapistServiceRes

	for _, service := range service.Physiotherapist.ServiceAndSchedule {
		if service.Id == serviceObjID {
			serviceRes := services.PhysiotherapistServiceRes{
				Id:          service.Id,
				Name:        service.Name,
				ServiceFees: service.ServiceFees,
			}

			if len(service.Slots) > 0 {
				for _, schedule := range service.Slots {
					serviceRes.Slots = append(serviceRes.Slots, services.Slots{
						StartTime: schedule.StartTime,
						EndTime:   schedule.EndTime,
						Days:      schedule.Days,
					})
				}
			}

			servicesRes = serviceRes
			break
		}
	}

	return c.Status(fiber.StatusOK).JSON(services.GetPhysiotherapistServicesResDto{
		Status:  true,
		Message: "Services retrieved successfully",
		Data:    servicesRes,
	})
}
