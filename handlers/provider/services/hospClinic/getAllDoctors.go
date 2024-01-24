package hospClinic

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary GetAllDoctors
// @Description GetAllDoctors
// @Tags hospClinic
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} services.DoctorResDto
// @Router /provider/services/get-all-doctors [get]
func GetAllDoctors(c *fiber.Ctx) error {
	ctx := context.TODO()

	var service entity.ServiceEntity

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"_id": providerData.ProviderId,
	}

	projection := bson.M{
		"hospClinic.doctor.id":         1,
		"hospClinic.doctor.name":       1,
		"hospClinic.doctor.speciality": 1,
		"hospClinic.doctor.image":      1,
		"hospClinic.doctor.schedule": bson.M{
			"startTime": 1,
			"endTime":   1,
			"days":      1,
		},
	}

	findOptions := options.FindOne().SetProjection(projection)

	err := serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.DoctorResDto{
				Status:  false,
				Message: "service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.DoctorResDto{
			Status:  false,
			Message: "Failed to fetch service from MongoDB: " + err.Error(),
		})
	}

	if service.HospClinic != nil {
		return c.Status(fiber.StatusNotFound).JSON(services.DoctorResDto{
			Status:  false,
			Message: "No HospClinic information found for the service",
		})
	}

	// Group doctors by speciality
	doctorsBySpeciality := make(map[string][]services.DoctorRes)

	// var doctorsRes []services.DoctorRes
	if service.HospClinic != nil && len(service.HospClinic.Doctor) > 0 {

		for _, doctor := range service.HospClinic.Doctor {
			doctorRes := services.DoctorRes{
				Id:         doctor.Id,
				Name:       doctor.Name,
				Image:      doctor.Image,
				Speciality: doctor.Speciality,
			}

			if len(doctor.Schedule) > 0 {
				for _, schedule := range doctor.Schedule {
					doctorRes.Schedule = append(doctorRes.Schedule, services.DoctorScheduleRes{
						StartTime: schedule.StartTime,
						EndTime:   schedule.EndTime,
						Days:      schedule.Days,
					})
				}
			}

			doctorsBySpeciality[doctor.Speciality] = append(doctorsBySpeciality[doctor.Speciality], doctorRes)
		}
	}
	// Convert the map to an array for response
	var response []services.SpecialityDoctorsRes

	for speciality, doctors := range doctorsBySpeciality {
		response = append(response, services.SpecialityDoctorsRes{
			Speciality: speciality,
			Doctors:    doctors,
		})
	}

	return c.Status(fiber.StatusOK).JSON(services.DoctorResDto{
		Status:  true,
		Message: "doctors retrieved successfully",
		Data:    response,
	})
}
