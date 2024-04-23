package hospitals

import (
	"careville_backend/database"
	hospitals "careville_backend/dto/customer/hospitals"
	"careville_backend/entity"
	helper "careville_backend/utils/helperFunctions"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Get all doctors
// @Description Retrieves information about doctors for a given service
// @Tags customer hospitals
// @Accept json
// @Produce json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param serviceId query string true "ID of the service"
// @Success 200 {object} DoctorResDto "Success response"
// @Failure 400 {object} DoctorResDto "Bad request"
// @Failure 404 {object} DoctorResDto "Not found"
// @Failure 500 {object} DoctorResDto "Internal server error"
// @Router /customer/healthFacility/get-all-doctors [get]
func GetAllDoctors(c *fiber.Ctx) error {

	serviceId := c.Query("serviceId")
	if serviceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(hospitals.DoctorResDto{
			Status:  false,
			Message: "service Id is mandatory",
		})
	}

	serviceObjectID, err := primitive.ObjectIDFromHex(serviceId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(hospitals.DoctorResDto{
			Status:  false,
			Message: "Invalid ID format",
		})
	}

	serviceColl := database.GetCollection("service")

	filter := bson.M{"_id": serviceObjectID}

	var service entity.ServiceEntity
	err = serviceColl.FindOne(ctx, filter).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(hospitals.DoctorResDto{
				Status:  false,
				Message: "Service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(hospitals.DoctorResDto{
			Status:  false,
			Message: "Failed to fetch service from MongoDB: " + err.Error(),
		})
	}

	if service.HospClinic == nil {
		return c.Status(fiber.StatusNotFound).JSON(hospitals.DoctorResDto{
			Status:  false,
			Message: "No HospClinic information found for the service",
		})
	}

	doctorsBySpeciality := make(map[string][]hospitals.DoctorRes)
	// Initialize the doctorsBySpeciality map with an empty slice for each speciality
	for _, doctor := range service.HospClinic.Doctor {
		doctorsBySpeciality[doctor.Speciality] = []hospitals.DoctorRes{}
	}
	if service.HospClinic != nil && len(service.HospClinic.Doctor) >= 1 {
		for _, doctor := range service.HospClinic.Doctor {
			nextAvailableSlots, _, err := GetDoctorNextAvailableDayAndSlots(doctor.Schedule)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(hospitals.DoctorResDto{
					Status:  false,
					Message: "Failed to get next available time slots",
				})
			}
			doctorRes := hospitals.DoctorRes{
				Id:            doctor.Id,
				Name:          doctor.Name,
				Image:         doctor.Image,
				Speciality:    doctor.Speciality,
				NextAvailable: nextAvailableSlots,
			}

			doctorsBySpeciality[doctor.Speciality] = append(doctorsBySpeciality[doctor.Speciality], doctorRes)
		}
	} else {

		if len(service.HospClinic.Doctor) == 1 {
			doctor := service.HospClinic.Doctor[0]
			doctorRes := hospitals.DoctorRes{
				Id:         doctor.Id,
				Name:       doctor.Name,
				Image:      doctor.Image,
				Speciality: doctor.Speciality,
				// NextAvailable: nextAvailable,
			}

			doctorsBySpeciality[doctor.Speciality] = append(doctorsBySpeciality[doctor.Speciality], doctorRes)
		}
	}

	var response []hospitals.SpecialityDoctorsRes

	for speciality, doctors := range doctorsBySpeciality {
		response = append(response, hospitals.SpecialityDoctorsRes{
			Speciality: speciality,
			Doctors:    doctors,
		})
	}

	return c.Status(fiber.StatusOK).JSON(hospitals.DoctorResDto{
		Status:  true,
		Message: "Doctors retrieved successfully",
		Data:    response,
	})
}

func GetDoctorNextAvailableDayAndSlots(schedules []entity.Schedule) (hospitals.NextAvailable, []entity.Schedule, error) {
	currentTime := time.Now().UTC()
	var nextAvailable hospitals.NextAvailable

	for _, schedule := range schedules {
		if !HasBreakingSlots(schedule.BreakingSlots) {
			for _, slot := range schedule.BreakingSlots {
				if helper.ContainsDay(schedule.Days, currentTime.Weekday().String()) && helper.DayAfterCurrentDay(schedule.Days[0], currentTime) {
					continue
				}
				for _, day := range schedule.Days {
					if helper.DayAfterCurrentDay(day, currentTime) {
						nextAvailable.StartTime = slot.StartTime
						return nextAvailable, []entity.Schedule{schedule}, nil
					}
				}
			}
		} else {
			nextAvailable.StartTime = GetUpcomingStartAndLastTime(schedules)
			if nextAvailable.StartTime != "" {
				return nextAvailable, []entity.Schedule{schedule}, nil
			}
		}
	}

	return nextAvailable, nil, errors.New("no next available slot found")
}

func HasBreakingSlots(slots []entity.BreakingSlots) bool {
	for _, slot := range slots {
		startTime, _ := time.Parse("15:04", slot.StartTime)
		endTime, _ := time.Parse("15:04", slot.EndTime)
		currentTime := time.Now().UTC()

		if currentTime.After(startTime) && currentTime.Before(endTime) {
			return true
		}
	}
	return false
}

func GetUpcomingStartAndLastTime(slots []entity.Schedule) string {
	var upcomingStartTime string
	currentTime := time.Now().UTC()

	for _, slot := range slots {
		// Check if the slot is for the current day or later
		if helper.ContainsDay(slot.Days, currentTime.Weekday().String()) && helper.DayAfterCurrentDay(slot.Days[0], currentTime) {
			// Check if the slot's start time is after the current time
			slotStartTime, err := time.Parse("15:04", slot.StartTime)
			if err != nil {
				continue // Skip this slot if parsing fails
			}
			if slotStartTime.After(currentTime) {
				// Update upcomingStartTime if it's empty or later than the current slot's start time
				if upcomingStartTime == "" || slot.StartTime < upcomingStartTime {
					upcomingStartTime = slot.StartTime
				}
			}
		}
	}

	return upcomingStartTime
}
