package doctorProfession

import (
	"careville_backend/database"
	"careville_backend/dto/customer/doctorProfession"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get doctorProfession by ID
// @Tags customer doctorProfession
// @Description Get doctorProfession by ID
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "doctorProfession ID"
// @Produce json
// @Success 200 {object} doctorProfession.GetDoctorProfessionResDto
// @Router /customer/healthProfessional/get-doctor/{id} [get]
func GetDoctorProfessionByID(c *fiber.Ctx) error {
	var (
		serviceColl = database.GetCollection("service")
	)

	idParam := c.Params("id")
	doctorProfessionID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(doctorProfession.GetDoctorProfessionResDto{
			Status:  false,
			Message: "Invalid doctorProfession ID",
		})
	}

	filter := bson.M{"_id": doctorProfessionID}

	projection := bson.M{
		"doctor.information.name":              1,
		"doctor.information.image":             1,
		"_id":                                  1,
		"doctor.information.additionalText":    1,
		"doctor.review.totalReviews":           1,
		"doctor.review.avgRating":              1,
		"doctor.additionalServices.speciality": 1,
		"doctor.schedule.consultationFees":     1,
		"doctor.schedule.slots": bson.M{
			"id":        1,
			"startTime": 1,
			"endTime":   1,
			"days":      1,
		},
	}

	findOptions := options.FindOne().SetProjection(projection)

	var doctorProfessionData entity.ServiceEntity
	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&doctorProfessionData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(doctorProfession.GetDoctorProfessionResDto{
			Status:  false,
			Message: "Failed to fetch doctorProfession data: " + err.Error(),
		})
	}

	if doctorProfessionData.Doctor == nil {
		return c.Status(fiber.StatusNotFound).JSON(doctorProfession.GetDoctorProfessionResDto{
			Status:  false,
			Message: "doctorProfession data not found",
		})
	}

	scheduleData := make([]doctorProfession.DoctorSchedule, 0)
	if doctorProfessionData.Doctor != nil && len(doctorProfessionData.Doctor.Schedule.Slots) > 0 {
		for _, service := range doctorProfessionData.Doctor.Schedule.Slots {
			scheduleData = append(scheduleData, doctorProfession.DoctorSchedule{
				Id:        service.Id,
				StartTime: service.StartTime,
				EndTime:   service.EndTime,
				Days:      service.Days,
			})
		}
	}

	doctorProfessionRes := doctorProfession.GetDoctorProfessionResDto{
		Status:  true,
		Message: "Doctor data fetched successfully",
		Data: doctorProfession.DoctorProfessionResponse{
			Id:               doctorProfessionData.Id,
			Image:            doctorProfessionData.Doctor.Information.Image,
			Name:             doctorProfessionData.Doctor.Information.Name,
			Speciality:       doctorProfessionData.Doctor.AdditionalServices.Speciality,
			AboutMe:          doctorProfessionData.Doctor.Information.AdditionalText,
			ConsultationFees: doctorProfessionData.Doctor.Schedule.ConsultationFees,
			DoctorSchedule:   scheduleData,
			TotalReviews:     doctorProfessionData.Doctor.Review.TotalReviews,
			AvgRating:        doctorProfessionData.Doctor.Review.AvgRating,
		},
	}

	return c.Status(fiber.StatusOK).JSON(doctorProfessionRes)
}
