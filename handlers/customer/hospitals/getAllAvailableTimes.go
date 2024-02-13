package hospitals

import (
	"careville_backend/database"
	hospitals "careville_backend/dto/customer/hospitals"
	"careville_backend/entity"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Get all available times
// @Description Get all available times
// @Tags customer hospitals
// @Accept application/json
// @Param hospitalId query string true "hospital ID"
// @Param doctorId query string true "doctor ID"
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} AvailableSlotsResDto
// @Router /customer/healthFacility/get-all-available-times [get]
func GetAllAvailableTimes(c *fiber.Ctx) error {

	var service entity.ServiceEntity
	var response []hospitals.Slots // Define response slice

	serviceColl := database.GetCollection("service")

	idParam := c.Query("hospitalId")
	hospitalId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(hospitals.AvailableSlotsResDto{
			Status:  false,
			Message: "Invalid hospital ID",
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

	filter := bson.M{
		"_id": hospitalId,
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

	currentDay := time.Now().Weekday()
	for _, doc := range service.HospClinic.Doctor {
		for _, sch := range doc.Schedule {
			var days []time.Weekday
			for _, day := range sch.Days {
				weekday, err := stringToWeekday(day)
				if err != nil {
					return c.Status(fiber.StatusBadRequest).JSON(hospitals.AvailableSlotsResDto{
						Status:  false,
						Message: "Invalid day: " + err.Error(),
					})
				}
				days = append(days, weekday)
			}
	
			if contains(days, currentDay) {
				startTime, err := time.Parse("15:04", sch.StartTime)
				if err != nil {
					return c.Status(400).JSON(hospitals.AvailableSlotsResDto{
						Status:  false,
						Message: "error parsing start time" + err.Error(),
					})
				}
	
				var filteredEvents []entity.UpcommingEvents
				for _, event := range doc.UpcommingEvents {
					if event.StartTime.Weekday() == currentDay && event.StartTime.After(time.Now().UTC()) {
						filteredEvents = append(filteredEvents, event)
					}
				}
	
				for i := 0; i < len(filteredEvents); i++ {
					if i == 0 {
						response = append(response, entityToHospitalsSlots(generateBreakingSlots(startTime.Format("15:04"), filteredEvents[i].StartTime.Format("15:04"), filteredEvents))...)
					} else {
						response = append(response, entityToHospitalsSlots(generateBreakingSlots(filteredEvents[i-1].EndTime.Format("15:04"), filteredEvents[i].StartTime.Format("15:04"), filteredEvents))...)
					}
				}
			}
		}
	}	

	if len(response) == 0 {
		return c.Status(fiber.StatusOK).JSON(hospitals.AvailableSlotsResDto{
			Status:  true,
			Message: "No slots available",
		})
	}

	return c.Status(fiber.StatusOK).JSON(hospitals.AvailableSlotsResDto{
		Status:  true,
		Message: "Available Slots retrieved successfully",
		Data:    response,
	})
}

func contains(days []time.Weekday, target time.Weekday) bool {
	for _, day := range days {
		if day == target {
			return true
		}
	}
	return false
}

func generateBreakingSlots(startTime, endTime string, upcomingEvents []entity.UpcommingEvents) []entity.Slots {
	var breakingSlots []entity.Slots
	layout := "15:04"

	start, _ := time.Parse(layout, startTime)
	end, _ := time.Parse(layout, endTime)

	for start.Before(end) {
		next := start.Add(30 * time.Minute)
		if next.After(end) {
			next = end
		}

		overlap := false
		for _, event := range upcomingEvents {
			if !(start.After(event.EndTime) || next.Before(event.StartTime)) {
				overlap = true
				break
			}
		}

		if !overlap {
			breakingSlots = append(breakingSlots, entity.Slots{
				StartTime: start.Format(layout),
				EndTime:   next.Format(layout),
			})
		}

		start = next
	}

	lastSlotEndTime, _ := time.Parse(layout, breakingSlots[len(breakingSlots)-1].EndTime)
	if lastSlotEndTime.Before(end) {
		breakingSlots = append(breakingSlots, entity.Slots{
			StartTime: lastSlotEndTime.Format(layout),
			EndTime:   end.Format(layout),
		})
	}

	return breakingSlots
}



func entityToHospitalsSlots(slots []entity.Slots) []hospitals.Slots {
	var converted []hospitals.Slots
	for _, s := range slots {
		converted = append(converted, hospitals.Slots{
			StartTime: s.StartTime,
			EndTime:   s.EndTime,
		})
	}
	return converted
}

func stringToWeekday(day string) (time.Weekday, error) {
	switch day {
	case "Sunday":
		return time.Sunday, nil
	case "Monday":
		return time.Monday, nil
	case "Tuesday":
		return time.Tuesday, nil
	case "Wednesday":
		return time.Wednesday, nil
	case "Thursday":
		return time.Thursday, nil
	case "Friday":
		return time.Friday, nil
	case "Saturday":
		return time.Saturday, nil
	default:
		return time.Sunday, fmt.Errorf("invalid day: %s", day)
	}
}
