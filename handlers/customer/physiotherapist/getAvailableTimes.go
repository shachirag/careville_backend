package physiotherapist

import (
	"careville_backend/database"
	physiotherapist "careville_backend/dto/customer/physiotherapist"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get all physiotherapist
// @Description Get all physiotherapist
// @Tags customer physiotherapist
// @Accept application/json
// @Param physiotherapistId query string true "service ID"
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} physiotherapist.AvailableSlotsResDto
// @Router /customer/healthProfessional/get-physiotherapist-available-slots [get]
func GetPhysiotherapistAvailableTimes(c *fiber.Ctx) error {

	var service entity.ServiceEntity

	serviceColl := database.GetCollection("service")

	serviceId := c.Query("physiotherapistId")
	if serviceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(physiotherapist.AvailableSlotsResDto{
			Status:  false,
			Message: "physiotheristId Id is mandatory",
		})
	}

	serviceObjectID, err := primitive.ObjectIDFromHex(serviceId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(physiotherapist.AvailableSlotsResDto{
			Status:  false,
			Message: "Invalid ID format",
		})
	}

	filter := bson.M{
		"_id": serviceObjectID,
	}

	projection := bson.M{
		"physiotherapist.upcommingEvents": 1,
		"physiotherapist.serviceAndSchedule.slots": bson.M{
			"startTime":     1,
			"endTime":       1,
			"days":          1,
			"breakingSlots": 1,
		},
	}

	findOptions := options.FindOne().SetProjection(projection)

	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(physiotherapist.AvailableSlotsResDto{
				Status:  false,
				Message: "service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(physiotherapist.AvailableSlotsResDto{
			Status:  false,
			Message: "Failed to fetch physiotherapist from MongoDB: " + err.Error(),
		})
	}

	var scheduleData []physiotherapist.Schedule
	var upcommingEvents []physiotherapist.UpcommingEvents

	if service.Physiotherapist != nil {
		for _, service := range service.Physiotherapist.ServiceAndSchedule {
			var breakinSlots []physiotherapist.BreakinSlots
			for _, slot := range service.Slots {
				for _, breakingSlot := range slot.BreakingSlots {
					breakinSlots = append(breakinSlots, physiotherapist.BreakinSlots{
						StartTime: breakingSlot.StartTime,
						EndTime:   breakingSlot.EndTime,
					})
				}
				scheduleData = append(scheduleData, physiotherapist.Schedule{
					StartTime:    slot.StartTime,
					EndTime:      slot.EndTime,
					Days:         slot.Days,
					BreakinSlots: breakinSlots,
				})
			}
		}

		for _, event := range service.Physiotherapist.UpcommingEvents {
			upcommingEvents = append(upcommingEvents, physiotherapist.UpcommingEvents{
				Id:        event.Id,
				StartTime: event.StartTime,
				EndTime:   event.EndTime,
			})
		}
	}

	if upcommingEvents == nil {
		upcommingEvents = make([]physiotherapist.UpcommingEvents, 0)
	}

	response := physiotherapist.AvailableSlotsResDto{
		Status: true,
		Data: physiotherapist.AvailableSlotsRes{
			Schedule:        scheduleData,
			UpcommingEvents: upcommingEvents,
		},
		Message: "Available times Retrieved successfully",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
