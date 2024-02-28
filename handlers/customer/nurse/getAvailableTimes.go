package nurse

import (
	"careville_backend/database"
	nurse "careville_backend/dto/customer/nurse"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get all nurse
// @Description Get all investigations
// @Tags customer nurse
// @Accept application/json
// @Param nurseId query string true "service ID"
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} nurse.AvailableSlotsResDto
// @Router /customer/healthProfessional/get-nurse-available-slots [get]
func GetNurseAvailableTimes(c *fiber.Ctx) error {

	var service entity.ServiceEntity

	serviceColl := database.GetCollection("service")

	serviceId := c.Query("nurseId")
	if serviceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(nurse.AvailableSlotsResDto{
			Status:  false,
			Message: "nurse Id is mandatory",
		})
	}

	serviceObjectID, err := primitive.ObjectIDFromHex(serviceId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(nurse.AvailableSlotsResDto{
			Status:  false,
			Message: "Invalid ID format",
		})
	}

	filter := bson.M{
		"_id": serviceObjectID,
	}

	projection := bson.M{
		"nurse.schedule.upcommingEvents": 1,
		"nurse.schedule.slots": bson.M{
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
			return c.Status(fiber.StatusNotFound).JSON(nurse.AvailableSlotsResDto{
				Status:  false,
				Message: "service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(nurse.AvailableSlotsResDto{
			Status:  false,
			Message: "Failed to fetch nurse from MongoDB: " + err.Error(),
		})
	}

	var scheduleData []nurse.Schedule
	var upcommingEvents []nurse.UpcommingEvents

	if service.Nurse != nil {
		for _, service := range service.Nurse.Schedule {
			var breakinSlots []nurse.BreakinSlots
			for _, slot := range service.Slots {
				for _, breakingSlot := range slot.BreakingSlots {
					breakinSlots = append(breakinSlots, nurse.BreakinSlots{
						StartTime: breakingSlot.StartTime,
						EndTime:   breakingSlot.EndTime,
					})
				}
				scheduleData = append(scheduleData, nurse.Schedule{
					StartTime:    slot.StartTime,
					EndTime:      slot.EndTime,
					Days:         slot.Days,
					BreakinSlots: breakinSlots,
				})
			}
		}

		for _, event := range service.Nurse.UpcommingEvents {
			upcommingEvents = append(upcommingEvents, nurse.UpcommingEvents{
				Id:        event.Id,
				StartTime: event.StartTime,
				EndTime:   event.EndTime,
			})
		}
	}

	if upcommingEvents == nil {
		upcommingEvents = make([]nurse.UpcommingEvents, 0)
	}

	response := nurse.AvailableSlotsResDto{
		Status: true,
		Data: nurse.AvailableSlotsRes{
			Schedule:        scheduleData,
			UpcommingEvents: upcommingEvents,
		},
		Message: "Available times Retrieved successfully",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
