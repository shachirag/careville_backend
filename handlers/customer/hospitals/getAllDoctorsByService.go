package hospitals

import (
	"careville_backend/database"
	hospitals "careville_backend/dto/customer/hospitals"
	"careville_backend/entity"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAllDoctors retrieves information about doctors for a given service.
// It returns a JSON response containing doctor details.
func GetAllDoctors(c *fiber.Ctx) error {
	ctx := context.Background()

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

	for _, doctor := range service.HospClinic.Doctor {
		var nextAvailable string

		for _, schedule := range doctor.Schedule {
			for _, breakingSlot := range schedule.BreakingSlots {
				startTime, err := time.Parse("15:04", breakingSlot.StartTime)
				if err != nil {
					return c.Status(fiber.StatusBadRequest).JSON(hospitals.DoctorResDto{
						Status:  false,
						Message: "Invalid start time format",
					})
				}

				// Check if the current time is before the start time of the breaking slot
				if time.Now().Before(startTime) {
					nextAvailable = startTime.Format("15:04")
					break
				}
			}

			if nextAvailable != "" {
				break
			}
		}

		if nextAvailable == "" {
			nextAvailable = "No slots available"
		}

		doctorRes := hospitals.DoctorRes{
			Id:            doctor.Id,
			Name:          doctor.Name,
			Image:         doctor.Image,
			Speciality:    doctor.Speciality,
			NextAvailable: nextAvailable,
		}

		doctorsBySpeciality[doctor.Speciality] = append(doctorsBySpeciality[doctor.Speciality], doctorRes)
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
