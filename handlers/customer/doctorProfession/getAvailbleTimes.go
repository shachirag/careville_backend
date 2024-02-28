package doctorProfession

import (
	"careville_backend/database"
	doctorProfession "careville_backend/dto/customer/doctorProfession"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get all doctorProfession
// @Description Get all doctorProfession
// @Tags customer doctorProfession
// @Accept application/json
// @Param doctorId query string true "service ID"
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} doctorProfession.AvailableSlotsResDto
// @Router /customer/healthProfessional/get-doctor-available-slots [get]
func GetDoctorAvailableTimes(c *fiber.Ctx) error {

	var service entity.ServiceEntity

	serviceColl := database.GetCollection("service")

	serviceId := c.Query("doctorId")
	if serviceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(doctorProfession.AvailableSlotsResDto{
			Status:  false,
			Message: "doctorId Id is mandatory",
		})
	}

	serviceObjectID, err := primitive.ObjectIDFromHex(serviceId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(doctorProfession.AvailableSlotsResDto{
			Status:  false,
			Message: "Invalid ID format",
		})
	}

	filter := bson.M{
		"_id": serviceObjectID,
	}

	projection := bson.M{
		"doctor.upcommingEvents": 1,
		"doctor.schedule.slots": bson.M{
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
			return c.Status(fiber.StatusNotFound).JSON(doctorProfession.AvailableSlotsResDto{
				Status:  false,
				Message: "service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(doctorProfession.AvailableSlotsResDto{
			Status:  false,
			Message: "Failed to fetch doctorProfession from MongoDB: " + err.Error(),
		})
	}

	var scheduleData []doctorProfession.Schedule
	var upcommingEvents []doctorProfession.UpcommingEvents

	if service.Doctor != nil {
		for _, slot := range service.Doctor.Schedule.Slots {
			var breakingSlots []doctorProfession.BreakinSlots
			for _, breakingSlot := range slot.BreakingSlots {
				breakingSlots = append(breakingSlots, doctorProfession.BreakinSlots{
					StartTime: breakingSlot.StartTime,
					EndTime:   breakingSlot.EndTime,
				})
			}
			scheduleData = append(scheduleData, doctorProfession.Schedule{
				StartTime:    slot.StartTime,
				EndTime:      slot.EndTime,
				Days:         slot.Days,
				BreakinSlots: breakingSlots,
			})
		}

		upcomingEvents := make([]doctorProfession.UpcommingEvents, len(service.Doctor.UpcommingEvents))
		for i, event := range service.Doctor.UpcommingEvents {
			upcomingEvents[i] = doctorProfession.UpcommingEvents{
				Id:        event.Id,
				StartTime: event.StartTime,
				EndTime:   event.EndTime,
			}
		}
	}

	if upcommingEvents == nil {
		upcommingEvents = make([]doctorProfession.UpcommingEvents, 0)
	}

	response := doctorProfession.AvailableSlotsResDto{
		Status: true,
		Data: doctorProfession.AvailableSlotsRes{
			Schedule:        scheduleData,
			UpcommingEvents: upcommingEvents,
		},
		Message: "Available times Retrieved successfully",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
