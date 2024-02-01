package doctorProfession

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

// @Summary Get slot info
// @Description Get slot info
// @Tags doctorProfession
// @Accept application/json
//
//	@Param Authorization header  string  true    "Authentication header"
//
// @Param slotId path string true "slotId"
// @Produce json
// @Success 200 {object} services.GetDoctorProfessionSlotsResDto
// @Router /provider/services/get-doctorProfession-slot-info/{slotId} [get]
func GetSlotInfo(c *fiber.Ctx) error {
	ctx := context.Background()

	var service entity.ServiceEntity

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	slotId := c.Params("slotId")
	slotObjID, err := primitive.ObjectIDFromHex(slotId)

	if err != nil {
		return c.Status(400).JSON(services.GetDoctorProfessionSlotsResDto{
			Status:  false,
			Message: "invalid objectId " + err.Error(),
		})
	}

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"_id": providerData.ProviderId,
		"doctor.schedule.slots": bson.M{
			"$elemMatch": bson.M{
				"id": slotObjID,
			},
		},
	}

	projection := bson.M{
		"doctor.schedule.slots": bson.M{
			"id":        1,
			"startTime": 1,
			"endTime":   1,
			"days":      1,
		},
	}

	findOptions := options.FindOne().SetProjection(projection)

	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.GetDoctorProfessionSlotsResDto{
				Status:  false,
				Message: "slot not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.GetDoctorProfessionSlotsResDto{
			Status:  false,
			Message: "Failed to fetch slot from MongoDB: " + err.Error(),
		})
	}

	if service.Doctor == nil {
		return c.Status(fiber.StatusNotFound).JSON(services.GetDoctorProfessionSlotsResDto{
			Status:  false,
			Message: "No slot information found for the service",
		})
	}

	var slotsRes services.GetDoctorProfessionSlotsResDto

	for _, slot := range service.Doctor.Schedule.Slots {
		if slot.Id == slotObjID {
			slotRes := services.DoctorProfessionSlotsResponseDto{
				Slots: []services.DoctorSlots{
					{
						Id:        slot.Id,
						StartTime: slot.StartTime,
						EndTime:   slot.EndTime,
						Days:      slot.Days,
					},
				},
			}
			slotsRes.Data = slotRes
			break
		}
	}

	if len(slotsRes.Data.Slots) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(services.GetDoctorProfessionSlotsResDto{
			Status:  false,
			Message: "Slot not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(services.GetDoctorProfessionSlotsResDto{
		Status:  true,
		Message: "Slot retrieved successfully",
		Data:    slotsRes.Data,
	})
}
