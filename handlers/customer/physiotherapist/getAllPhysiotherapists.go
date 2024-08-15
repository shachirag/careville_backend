package physiotherapist

import (
	"careville_backend/database"
	physiotherapist "careville_backend/dto/customer/physiotherapist"
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

// @Summary Get physiotherapist
// @Tags customer physiotherapist
// @Description Get physiotherapist
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param search query string false "Filter physiotherapist by search"
// @Param long query float64 false "Longitude for memories sorting (required for distance sorting)"
// @Param lat query float64 false "Latitude for memories sorting (required for distance sorting)"
// @Produce json
// @Success 200 {object} physiotherapist.GetPhysiotherapistResponseDto
// @Router /customer/healthProfessional/get-physiotherapists [get]
func GetPhysiotherapists(c *fiber.Ctx) error {

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
			return c.Status(fiber.StatusBadRequest).JSON(physiotherapist.GetPhysiotherapistResponseDto{
				Status:  false,
				Message: "Invalid latitude format",
			})
		}

		long, err = strconv.ParseFloat(longParam, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(physiotherapist.GetPhysiotherapistResponseDto{
				Status:  false,
				Message: "Invalid longitude format",
			})
		}
	}

	filter := bson.M{
		"role":                 "healthProfessional",
		"facilityOrProfession": "physiotherapist",
	}

	maxDistance := 50000

	if latParam != "" && longParam != "" {
		filter["physiotherapist.information.address"] = bson.M{
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
		filter["physiotherapist.information.name"] = bson.M{"$regex": searchTitle, "$options": "i"}
	}

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	cursor, err := serviceColl.Find(ctx, filter, sortOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(physiotherapist.GetPhysiotherapistResponseDto{
			Status:  false,
			Message: "Failed to fetch physiotherapist data: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var physiotherapistData []physiotherapist.GetPhysiotherapistRes
	for cursor.Next(ctx) {
		var physiotherapist1 entity.ServiceEntity
		if err := cursor.Decode(&physiotherapist1); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(physiotherapist.GetPhysiotherapistResponseDto{
				Status:  false,
				Message: "Failed to decode physiotherapist data: " + err.Error(),
			})
		}

		if physiotherapist1.Physiotherapist != nil {
			// nextAvailableSlots, _, err := GetPhysiotherapistNextAvailableDayAndSlots(physiotherapist1.Physiotherapist.ServiceAndSchedule)
			// if err != nil {
			// 	return c.Status(fiber.StatusInternalServerError).JSON(nurse.GetNurseResponseDto{
			// 		Status:  false,
			// 		Message: "Failed to get next available time slots",
			// 	})
			// }
			physiotherapistData = append(physiotherapistData, physiotherapist.GetPhysiotherapistRes{
				Id:          physiotherapist1.Id,
				Image:       physiotherapist1.Physiotherapist.Information.Image,
				Name:        physiotherapist1.Physiotherapist.Information.Name,
				ServiceType: "Physiotherapist",
				NextAvailable: physiotherapist.NextAvailable{
					StartTime: "",
				},
			})
		}
	}

	if len(physiotherapistData) == 0 {
		return c.Status(fiber.StatusOK).JSON(physiotherapist.GetPhysiotherapistResponseDto{
			Status:  false,
			Message: "No Physiotherapist data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(physiotherapist.GetPhysiotherapistResponseDto{
		Status:  true,
		Message: "Successfully fetched physiotherapists data.",
		Data:    physiotherapistData,
	})
}

func GetPhysiotherapistNextAvailableDayAndSlots(schedules []entity.ServiceAndSchedule) (physiotherapist.NextAvailable, []entity.ServiceAndSchedule, error) {
	currentTime := time.Now().UTC()
	var nextAvailable physiotherapist.NextAvailable
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
