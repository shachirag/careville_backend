package medicalLabScientist

import (
	"careville_backend/database"
	medicalLabScientist "careville_backend/dto/customer/medicalLabScientist"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get all medicalLabScientist
// @Description Get all investigations
// @Tags customer medicalLabScientist
// @Accept application/json
// @Param medicalLabScientistId query string true "service ID"
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} medicalLabScientist.AvailableSlotsResDto
// @Router /customer/healthProfessional/get-medicalLabScientist-available-slots [get]
func GetMedicalLabScientistAvailableTimes(c *fiber.Ctx) error {

	var service entity.ServiceEntity

	serviceColl := database.GetCollection("service")

	serviceId := c.Query("medicalLabScientistId")
	if serviceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(medicalLabScientist.AvailableSlotsResDto{
			Status:  false,
			Message: "hospital Id is mandatory",
		})
	}

	serviceObjectID, err := primitive.ObjectIDFromHex(serviceId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(medicalLabScientist.AvailableSlotsResDto{
			Status:  false,
			Message: "Invalid ID format",
		})
	}

	filter := bson.M{
		"_id": serviceObjectID,
	}

	projection := bson.M{
		"medicalLabScientist.upcommingEvents": 1,
		"medicalLabScientist.serviceAndSchedule.slots": bson.M{
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
			return c.Status(fiber.StatusNotFound).JSON(medicalLabScientist.AvailableSlotsResDto{
				Status:  false,
				Message: "service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.AvailableSlotsResDto{
			Status:  false,
			Message: "Failed to fetch medicalLabScientist from MongoDB: " + err.Error(),
		})
	}

	var scheduleData []medicalLabScientist.Schedule
	var upcommingEvents []medicalLabScientist.UpcommingEvents

	if service.MedicalLabScientist != nil {
		for _, service := range service.MedicalLabScientist.ServiceAndSchedule {
			var breakinSlots []medicalLabScientist.BreakinSlots
			for _, slot := range service.Slots {
				for _, breakingSlot := range slot.BreakingSlots {
					breakinSlots = append(breakinSlots, medicalLabScientist.BreakinSlots{
						StartTime: breakingSlot.StartTime,
						EndTime:   breakingSlot.EndTime,
					})
				}
				scheduleData = append(scheduleData, medicalLabScientist.Schedule{
					StartTime:    slot.StartTime,
					EndTime:      slot.EndTime,
					Days:         slot.Days,
					BreakinSlots: breakinSlots,
				})
			}
		}

		for _, event := range service.MedicalLabScientist.UpcommingEvents {
			upcommingEvents = append(upcommingEvents, medicalLabScientist.UpcommingEvents{
				Id:        event.Id,
				StartTime: event.StartTime,
				EndTime:   event.EndTime,
			})
		}
	}

	if upcommingEvents == nil {
		upcommingEvents = make([]medicalLabScientist.UpcommingEvents, 0)
	}

	response := medicalLabScientist.AvailableSlotsResDto{
		Status: true,
		Data: medicalLabScientist.AvailableSlotsRes{
			Schedule:        scheduleData,
			UpcommingEvents: upcommingEvents,
		},
		Message: "medicalLabScientist retrieved successfully",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
