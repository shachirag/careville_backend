package medicalLabScientist

import (
	"careville_backend/database"
	"careville_backend/dto/customer/medicalLabScientist"
	"careville_backend/entity"
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

func GetMedicalLabScientist1(c *fiber.Ctx) error {
	var (
		serviceColl = database.GetCollection("service")
	)

	searchTitle := c.Query("search", "")

	var lat, long float64
	latParam := c.Query("lat")
	longParam := c.Query("long")
	var err error

	if latParam != "" && longParam != "" {
		lat, err = strconv.ParseFloat(latParam, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(medicalLabScientist.GetMedicalLabScientistResponseDto{
				Status:  false,
				Message: "Invalid latitude format",
			})
		}

		long, err = strconv.ParseFloat(longParam, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(medicalLabScientist.GetMedicalLabScientistResponseDto{
				Status:  false,
				Message: "Invalid longitude format",
			})
		}
	}

	filter := bson.M{
		"role":                 "healthProfessional",
		"facilityOrProfession": "medicalLabScientist",
	}

	maxDistance := 50000

	if latParam != "" && longParam != "" {
		filter["medicalLabScientist.information.address"] = bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{long, lat},
				},
				"$maxDistance": maxDistance,
			},
		}
	}

	if searchTitle != "" {
		filter["medicalLabScientist.information.name"] = bson.M{"$regex": searchTitle, "$options": "i"}
	}

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	projection := bson.M{
		"medicalLabScientist.information.name":  1,
		"medicalLabScientist.information.image": 1,
		"_id":                                   1,
	}

	findOptions := options.Find().SetProjection(projection)

	cursor, err := serviceColl.Find(ctx, filter, sortOptions, findOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.GetMedicalLabScientistResponseDto{
			Status:  false,
			Message: "Failed to fetch medicalLabScientist data: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var medicalLabScientistData []medicalLabScientist.GetMedicalLabScientistRes
	for cursor.Next(ctx) {
		var medicalLabScientist1 entity.ServiceEntity
		if err := cursor.Decode(&medicalLabScientist1); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.GetMedicalLabScientistResponseDto{
				Status:  false,
				Message: "Failed to decode medicalLabScientist data: " + err.Error(),
			})
		}
		if medicalLabScientist1.MedicalLabScientist != nil {
			nextAvailableSlots, err := getNextAvailableDayAndSlots(medicalLabScientist1.MedicalLabScientist.ServiceAndSchedule)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.GetMedicalLabScientistResponseDto{
					Status:  false,
					Message: "Failed to get next available time slots: ",
				})
			}
			medicalLabScientistData = append(medicalLabScientistData, medicalLabScientist.GetMedicalLabScientistRes{
				Id:            medicalLabScientist1.Id,
				Image:         medicalLabScientist1.MedicalLabScientist.Information.Image,
				Name:          medicalLabScientist1.MedicalLabScientist.Information.Name,
				ServiceType:   "MedicalLabScientist",
				NextAvailable: nextAvailableSlots,
			})
		}
	}

	if len(medicalLabScientistData) == 0 {
		return c.Status(fiber.StatusOK).JSON(medicalLabScientist.GetMedicalLabScientistResponseDto{
			Status:  false,
			Message: "No MedicalLabScientist data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(medicalLabScientist.GetMedicalLabScientistResponseDto{
		Status:  true,
		Message: "Successfully fetched medicalLabScientists data.",
		Data:    medicalLabScientistData,
	})
}

func getNextAvailableDayAndSlots(schedules []entity.ServiceAndSchedule) (medicalLabScientist.NextAvailable, []entity.ServiceAndSchedule) {
    currentTime := time.Now()
    var nextAvailable medicalLabScientist.NextAvailable
    for _, schedule := range schedules {
        if !hasBreakingSlots(schedule.Slots) {
            for _, slot := range schedule.Slots {
                if containsDay(slot.Days, currentTime.Weekday().String()) && dayAfterCurrentDay(slot.Days[0], currentTime) {
                    nextAvailable.StartTime = slot.Days[0]
                    nextAvailable.LastTime = getLastTimeAvailable(schedule.Slots)
                    return nextAvailable, []entity.ServiceAndSchedule{schedule}
                }
            }
        } else {
            nextAvailable.StartTime, nextAvailable.LastTime = getUpcomingStartAndLastTime(schedule.Slots)
            if nextAvailable.StartTime != "" && nextAvailable.LastTime != "" {
                return nextAvailable, []entity.ServiceAndSchedule{schedule}
            }
        }
    }
    return medicalLabScientist.NextAvailable{}, nil
}

func getUpcomingStartAndLastTime(slots []entity.Slots) (string, string) {
    currentTime := time.Now()
    for _, slot := range slots {
        for _, breakingSlot := range slot.BreakingSlots {
            startTime, _ := time.Parse("15:04", breakingSlot.StartTime)
            endTime, _ := time.Parse("15:04", breakingSlot.EndTime)
            if startTime.After(currentTime) && endTime.After(startTime) {
                return breakingSlot.StartTime, breakingSlot.EndTime
            }
        }
    }
    return "", ""
}

func getLastTimeAvailable(slots []entity.Slots) string {
    var lastEndTime string
    for _, slot := range slots {
        for _, breakingSlot := range slot.BreakingSlots {
            endTime := breakingSlot.EndTime
            if endTime > lastEndTime {
                lastEndTime = endTime
            }
        }
    }
    return lastEndTime
}

func hasBreakingSlots(slots []entity.Slots) bool {
	for _, slot := range slots {
		for _, breakingSlot := range slot.BreakingSlots {
			startTime, _ := time.Parse("15:04", breakingSlot.StartTime)
			endTime, _ := time.Parse("15:04", breakingSlot.EndTime)
			currentTime := time.Now()

			if currentTime.After(startTime) && currentTime.Before(endTime) {
				return true
			}
		}
	}
	return false
}

func containsDay(days []string, target string) bool {
	for _, day := range days {
		if day == target {
			return true
		}
	}
	return false
}

func dayAfterCurrentDay(day string, currentTime time.Time) bool {
	currentWeekday := currentTime.Weekday().String()
	if day == currentWeekday {
		return false
	}

	daysMap := map[string]int{
		"Sunday":    0,
		"Monday":    1,
		"Tuesday":   2,
		"Wednesday": 3,
		"Thursday":  4,
		"Friday":    5,
		"Saturday":  6,
	}

	currentDayNum := daysMap[currentWeekday]
	targetDayNum := daysMap[day]

	return currentDayNum < targetDayNum
}
