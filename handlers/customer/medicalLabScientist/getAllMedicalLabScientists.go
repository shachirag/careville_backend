package medicalLabScientist

import (
	"careville_backend/database"
	"careville_backend/dto/customer/medicalLabScientist"
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

	cursor, err := serviceColl.Find(ctx, filter, sortOptions)
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
			// nextAvailableSlots, _, err := GetNextAvailableDayAndSlots(medicalLabScientist1.MedicalLabScientist.ServiceAndSchedule)
			// if err != nil {
			// 	return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.GetMedicalLabScientistResponseDto{
			// 		Status:  false,
			// 		Message: "Failed to get next available time slots: " + err.Error(),
			// 	})
			// }

			medicalLabScientistData = append(medicalLabScientistData, medicalLabScientist.GetMedicalLabScientistRes{
				Id:          medicalLabScientist1.Id,
				Image:       medicalLabScientist1.MedicalLabScientist.Information.Image,
				Name:        medicalLabScientist1.MedicalLabScientist.Information.Name,
				ServiceType: "MedicalLabScientist",
				NextAvailable: medicalLabScientist.NextAvailable{
					StartTime: "",
				},
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

func GetNextAvailableDayAndSlots(schedules []entity.ServiceAndSchedule) (medicalLabScientist.NextAvailable, []entity.ServiceAndSchedule, error) {
	currentTime := time.Now().UTC()
	var nextAvailable medicalLabScientist.NextAvailable
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
