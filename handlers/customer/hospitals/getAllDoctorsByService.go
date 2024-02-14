package hospitals

import (
	"careville_backend/database"
	hospitals "careville_backend/dto/customer/hospitals"
	"careville_backend/entity"
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
// @Param serviceId query string true "ID of the service"
// @Success 200 {object} DoctorResDto "Success response"
// @Failure 400 {object} DoctorResDto "Bad request"
// @Failure 404 {object} DoctorResDto "Not found"
// @Failure 500 {object} DoctorResDto "Internal server error"
// @Router /customer/hospitals/get-all-doctors [get]
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

	if service.HospClinic != nil && len(service.HospClinic.Doctor) > 1 {
		for _, doctor := range service.HospClinic.Doctor {
			nextAvailable := hospitals.NextAvailable{}

			for _, schedule := range doctor.Schedule {
				for _, breakingSlot := range schedule.BreakingSlots {
					startTime, err := time.Parse("15:04", breakingSlot.StartTime)
					if err != nil {
						return c.Status(fiber.StatusBadRequest).JSON(hospitals.DoctorResDto{
							Status:  false,
							Message: "Invalid start time format",
						})
					}

					startTimeUTC := startTime.UTC()

					if startTimeUTC.After(time.Now()) {
						nextAvailable.StartTime = startTimeUTC.Format("15:04")
						nextAvailable.LastTime = breakingSlot.EndTime
						break
					}

				}

				// if !nextAvailable.IsEmpty() {
				// 	break
				// }
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
