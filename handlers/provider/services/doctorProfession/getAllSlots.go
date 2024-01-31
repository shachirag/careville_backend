package doctorProfession

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

// @Summary Get all slots
// @Description Get all slots
// @Tags doctorProfession
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} services.DoctorProfessionSlotsResDto
// @Router /provider/services/get-doctorProfession-slots [get]
func GetDoctorProfessionSlots(c *fiber.Ctx) error {
	ctx := context.Background()

	var service entity.ServiceEntity

	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"_id": providerData.ProviderId,
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

	err := serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.DoctorProfessionSlotsResDto{
				Status:  false,
				Message: "slot not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.DoctorProfessionSlotsResDto{
			Status:  false,
			Message: "Failed to fetch slot from MongoDB: " + err.Error(),
		})
	}

	// Check if the slots are nil or empty
	if service.Doctor.Schedule.Slots == nil || len(service.Doctor.Schedule.Slots) == 0 {
		return c.Status(fiber.StatusOK).JSON(services.DoctorProfessionSlotsResDto{
			Status:  false,
			Message: "No slot data found.",
		})
	}

	// Adjusting the loop to iterate over DoctorSlots directly
	var slotData []services.DoctorSlots
	for _, slot := range service.Doctor.Schedule.Slots {
		slotData = append(slotData, services.DoctorSlots{
			Id:        slot.Id,
			StartTime: slot.StartTime,
			EndTime:   slot.EndTime,
			Days:      slot.Days,
		})
	}

	// Wrap slotData inside DoctorProfessionSlotsResponseDto
	var responseDto services.DoctorProfessionSlotsResDto
	for _, slot := range slotData {
		responseDto.Data = append(responseDto.Data, services.DoctorProfessionSlotsResponseDto{
			Slots: []services.DoctorSlots{slot},
		})
	}

	if len(responseDto.Data) == 0 {
		return c.Status(fiber.StatusOK).JSON(services.DoctorProfessionSlotsResDto{
			Status:  false,
			Message: "No slot data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(responseDto)
}
