package medicalLabScientist

import (
	"careville_backend/database"
	medicalLabScientist "careville_backend/dto/customer/medicalLabScientist"
	"careville_backend/entity"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get medicalLabScientist by ID
// @Tags customer medicalLabScientist
// @Description Get medicalLabScientist by ID
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "medicalLabScientist ID"
// @Produce json
// @Success 200 {object} medicalLabScientist.GetMedicalLabScientistResDto
// @Router /customer/healthProfessional/get-medicalLabScientist/{id} [get]
func GetMedicalLabScientistByID(c *fiber.Ctx) error {
	var (
		serviceColl = database.GetCollection("service")
	)

	idParam := c.Params("id")
	medicalLabScientistID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(medicalLabScientist.GetMedicalLabScientistResDto{
			Status:  false,
			Message: "Invalid medicalLabScientist ID",
		})
	}

	filter := bson.M{"_id": medicalLabScientistID}

	projection := bson.M{
		"medicalLabScientist.information.name":               1,
		"medicalLabScientist.information.image":              1,
		"_id":                                                1,
		"medicalLabScientist.information.additionalText":     1,
		"medicalLabScientist.review.totalReviews":            1,
		"medicalLabScientist.review.avgRating":               1,
		"medicalLabScientist.serviceAndSchedule.id":          1,
		"medicalLabScientist.serviceAndSchedule.name":        1,
		"medicalLabScientist.serviceAndSchedule.serviceFees": 1,
		"medicalLabScientist.serviceAndSchedule.slots": bson.M{
			"startTime": 1,
			"endTime":   1,
			"days":      1,
		},
	}

	findOptions := options.FindOne().SetProjection(projection)
	fmt.Print("1233")
	var medicalLabScientistData entity.ServiceEntity
	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&medicalLabScientistData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.GetMedicalLabScientistResDto{
			Status:  false,
			Message: "Failed to fetch medicalLabScientist data: " + err.Error(),
		})
	}

	if medicalLabScientistData.MedicalLabScientist == nil {
		return c.Status(fiber.StatusNotFound).JSON(medicalLabScientist.GetMedicalLabScientistResDto{
			Status:  false,
			Message: "medicalLabScientist data not found",
		})
	}
	fmt.Print("123")
	scheduleData := make([]medicalLabScientist.ServiceAndSchedule, 0)
	if medicalLabScientistData.MedicalLabScientist != nil && len(medicalLabScientistData.MedicalLabScientist.ServiceAndSchedule) > 0 {
		for _, service := range medicalLabScientistData.MedicalLabScientist.ServiceAndSchedule {
			var slotData []medicalLabScientist.Slots
			for _, slots := range service.Slots {
				slotData = append(slotData, medicalLabScientist.Slots{
					StartTime: slots.StartTime,
					EndTime:   slots.EndTime,
					Days:      slots.Days,
				})
			}
			scheduleData = append(scheduleData, medicalLabScientist.ServiceAndSchedule{
				Id:          service.Id,
				ServiceFees: service.ServiceFees,
				Name:        service.Name,
				Slots:       slotData,
			})
		}

	}

	var avgRating float64
	var totalReviews int32
	if medicalLabScientistData.MedicalLabScientist != nil && medicalLabScientistData.MedicalLabScientist.Review != nil {
		avgRating = medicalLabScientistData.MedicalLabScientist.Review.AvgRating
		totalReviews = medicalLabScientistData.MedicalLabScientist.Review.TotalReviews
	}

	medicalLabScientistRes := medicalLabScientist.GetMedicalLabScientistResDto{
		Status:  true,
		Message: "MedicalLabScientist data fetched successfully",
		Data: medicalLabScientist.MedicalLabScientistResponse{
			Id:                 medicalLabScientistData.Id,
			Image:              medicalLabScientistData.MedicalLabScientist.Information.Image,
			Name:               medicalLabScientistData.MedicalLabScientist.Information.Name,
			AboutMe:            medicalLabScientistData.MedicalLabScientist.Information.AdditionalText,
			ServiceAndSchedule: scheduleData,
			TotalReviews:       totalReviews,
			AvgRating:          avgRating,
		},
	}

	return c.Status(fiber.StatusOK).JSON(medicalLabScientistRes)
}
