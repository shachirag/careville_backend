package hospitals

import (
	"careville_backend/database"
	"careville_backend/dto/customer/hospitals"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Get available times
// @Description Get available times
// @Tags customer hospitals
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param hospitalId query string true "service ID"
// @Param doctorId query string true "doctor ID"
// @Produce json
// @Success 200 {object} hospitals.AvailableSlotsResDto
// @Router /customer/healthFacility/get-all-available-slots [get]
func GetAvailableSlots(c *fiber.Ctx) error {
	var service entity.ServiceEntity

	serviceId := c.Query("hospitalId")
	if serviceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(hospitals.AvailableSlotsResDto{
			Status:  false,
			Message: "hospital Id is mandatory",
		})
	}

	serviceObjectID, err := primitive.ObjectIDFromHex(serviceId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(hospitals.AvailableSlotsResDto{
			Status:  false,
			Message: "Invalid ID format",
		})
	}

	doctorId := c.Query("doctorId")
	doctorObjID, err := primitive.ObjectIDFromHex(doctorId)
	if err != nil {
		return c.Status(400).JSON(hospitals.AvailableSlotsResDto{
			Status:  false,
			Message: "invalid objectId " + err.Error(),
		})
	}

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"_id": serviceObjectID,
		"hospClinic.doctor": bson.M{
			"$elemMatch": bson.M{
				"id": doctorObjID,
			},
		},
	}

	err = serviceColl.FindOne(ctx, filter).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(hospitals.AvailableSlotsResDto{
				Status:  false,
				Message: "Service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(hospitals.AvailableSlotsResDto{
			Status:  false,
			Message: "Failed to fetch service from MongoDB: " + err.Error(),
		})
	}

	if service.HospClinic == nil {
		return c.Status(fiber.StatusNotFound).JSON(hospitals.AvailableSlotsResDto{
			Status:  false,
			Message: "No HospClinic information found for the service",
		})
	}

	var scheduleRes []hospitals.Schedule
	var eventRes []hospitals.UpcommingEvents

	for _, doctor := range service.HospClinic.Doctor {
		if doctor.Id == doctorObjID {
			if len(doctor.Schedule) > 0 {
				for _, schedule := range doctor.Schedule {
					var breakinSlots []hospitals.BreakinSlots
					if len(schedule.BreakingSlots) > 0 {
						for _, breakinSlot := range schedule.BreakingSlots {
							breakinSlots = append(breakinSlots, hospitals.BreakinSlots{
								StartTime: breakinSlot.StartTime,
								EndTime:   breakinSlot.EndTime,
							})
						}
					}
					scheduleRes = append(scheduleRes, hospitals.Schedule{
						StartTime:    schedule.StartTime,
						EndTime:      schedule.EndTime,
						Days:         schedule.Days,
						BreakinSlots: breakinSlots,
					})
				}
			}

			if len(doctor.UpcommingEvents) > 0 {
				for _, events := range doctor.UpcommingEvents {
					eventRes = append(eventRes, hospitals.UpcommingEvents{
						Id:        events.Id,
						StartTime: events.StartTime,
						EndTime:   events.EndTime,
					})
				}
			}

			break
		}
	}

	return c.Status(fiber.StatusOK).JSON(hospitals.AvailableSlotsResDto{
		Status:  true,
		Message: "Doctor retrieved successfully",
		Data: hospitals.AvailableSlotsRes{
			Schedule:        scheduleRes,
			UpcommingEvents: eventRes,
		},
	})
}
