package nurse

import (
	"careville_backend/database"
	nurse "careville_backend/dto/customer/nurse"
	"careville_backend/entity"
	helper "careville_backend/utils/helperFunctions"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

// @Summary Get nurse
// @Tags customer nurse
// @Description Get nurse
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param search query string false "Filter nurse by search"
// @Param long query float64 false "Longitude for memories sorting (required for distance sorting)"
// @Param lat query float64 false "Latitude for memories sorting (required for distance sorting)"
// @Produce json
// @Success 200 {object} nurse.GetNurseResponseDto
// @Router /customer/healthProfessional/get-nurses [get]
func GetNurses(c *fiber.Ctx) error {

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
			return c.Status(fiber.StatusBadRequest).JSON(nurse.GetNurseResponseDto{
				Status:  false,
				Message: "Invalid latitude format",
			})
		}

		long, err = strconv.ParseFloat(longParam, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(nurse.GetNurseResponseDto{
				Status:  false,
				Message: "Invalid longitude format",
			})
		}
	}

	filter := bson.M{
		"role":                 "healthProfessional",
		"facilityOrProfession": "nurse",
	}

	maxDistance := 50000

	if latParam != "" && longParam != "" {
		filter["nurse.information.address"] = bson.M{
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
		filter["nurse.information.name"] = bson.M{"$regex": searchTitle, "$options": "i"}
	}

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	cursor, err := serviceColl.Find(ctx, filter, sortOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(nurse.GetNurseResponseDto{
			Status:  false,
			Message: "Failed to fetch nurse data: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var nurseData []nurse.GetNurseRes
	for cursor.Next(ctx) {
		var nurse1 entity.ServiceEntity
		if err := cursor.Decode(&nurse1); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(nurse.GetNurseResponseDto{
				Status:  false,
				Message: "Failed to decode nurse data: " + err.Error(),
			})
		}
		if nurse1.Nurse != nil {
			nextAvailableSlots, _, err := GetNurseNextAvailableDayAndSlots(nurse1.Nurse.Schedule)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(nurse.GetNurseResponseDto{
					Status:  false,
					Message: "Failed to get next available time slots: " + err.Error(),
				})
			}
			nurseData = append(nurseData, nurse.GetNurseRes{
				Id:            nurse1.Id,
				Image:         nurse1.Nurse.Information.Image,
				Name:          nurse1.Nurse.Information.Name,
				ServiceType:   "Nurse",
				NextAvailable: nextAvailableSlots,
			})
		}
	}

	if len(nurseData) == 0 {
		return c.Status(fiber.StatusOK).JSON(nurse.GetNurseResponseDto{
			Status:  false,
			Message: "No Nurse data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(nurse.GetNurseResponseDto{
		Status:  true,
		Message: "Successfully fetched nurses data.",
		Data:    nurseData,
	})
}

func GetNurseNextAvailableDayAndSlots(schedules []entity.ServiceAndSchedule) (nurse.NextAvailable, []entity.ServiceAndSchedule, error) {
	currentTime := time.Now().UTC()
	var nextAvailable nurse.NextAvailable
	// fmt.Println("Length of schedules:", len(schedules))
	for _, schedule := range schedules {
		if !helper.HasBreakingSlots(schedule.Slots) {

			for _, slot := range schedule.Slots {
				if helper.ContainsDay(slot.Days, currentTime.Weekday().String()) && helper.DayAfterCurrentDay(slot.Days[0], currentTime) {
					continue
				}
				for _, day := range slot.Days {
					if helper.DayAfterCurrentDay(day, currentTime) {
						nextAvailable.StartTime = slot.StartTime
						return nextAvailable, []entity.ServiceAndSchedule{schedule}, nil
					}
				}
			}
		} else {
			nextAvailable.StartTime = helper.GetUpcomingStartAndLastTime(schedule.Slots)
			if nextAvailable.StartTime != "" {
				return nextAvailable, []entity.ServiceAndSchedule{schedule}, nil
			}
		}
	}

	return nextAvailable, nil, errors.New("no next available slot found")
}
