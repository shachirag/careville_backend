package hospClinic

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get Doctor info
// @Description Get Doctor info
// @Tags hospClinic
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param doctorId path string true "doctor ID"
// @Produce json
// @Success 200 {object} services.DoctorResDto
// @Router /provider/services/get-doctor-info/{doctorId} [get]
func GetDoctorsInfo(c *fiber.Ctx) error {
	ctx := context.TODO()

	var service entity.ServiceEntity

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	doctorId := c.Params("doctorId")
	doctorObjID, err := primitive.ObjectIDFromHex(doctorId)

	if err != nil {
		return c.Status(400).JSON(services.GetDoctorResDto{
			Status:  false,
			Message: "invalid objectId " + err.Error(),
		})
	}

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"_id": providerData.ProviderId,
		"hospClinic.doctor": bson.M{
			"$elemMatch": bson.M{
				"id": doctorObjID,
			},
		},
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

	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.GetDoctorResDto{
				Status:  false,
				Message: "Service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.GetDoctorResDto{
			Status:  false,
			Message: "Failed to fetch service from MongoDB: " + err.Error(),
		})
	}

	if service.HospClinic == nil {
		return c.Status(fiber.StatusNotFound).JSON(services.GetDoctorResDto{
			Status:  false,
			Message: "No HospClinic information found for the service",
		})
	}

	var doctorsRes services.DoctorRes

	for _, doctor := range service.HospClinic.Doctor {
		if doctor.Id == doctorObjID {
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

			doctorsRes = doctorRes
			break
		}
	}

	return c.Status(fiber.StatusOK).JSON(services.GetDoctorResDto{
		Status:  true,
		Message: "Doctor retrieved successfully",
		Data:    doctorsRes,
	})
}
