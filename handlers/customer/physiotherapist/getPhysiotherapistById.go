package physiotherapist

import (
	"careville_backend/database"
	physiotherapist "careville_backend/dto/customer/physiotherapist"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get physiotherapist by ID
// @Tags customer physiotherapist
// @Description Get physiotherapist by ID
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "physiotherapist ID"
// @Produce json
// @Success 200 {object} physiotherapist.GetPhysiotherapistResDto
// @Router /customer/healthProfessional/get-physiotherapist/{id} [get]
func GetPhysiotherapistByID(c *fiber.Ctx) error {
	var (
		serviceColl = database.GetCollection("service")
	)

	idParam := c.Params("id")
	physiotherapistID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(physiotherapist.GetPhysiotherapistResDto{
			Status:  false,
			Message: "Invalid physiotherapist ID",
		})
	}

	filter := bson.M{"_id": physiotherapistID}

	projection := bson.M{
		"physiotherapist.information.name":               1,
		"physiotherapist.information.image":              1,
		"_id":                                            1,
		"physiotherapist.information.additionalText":     1,
		"physiotherapist.review.totalReviews":            1,
		"physiotherapist.review.avgRating":               1,
		"physiotherapist.serviceAndSchedule.id":          1,
		"physiotherapist.serviceAndSchedule.name":        1,
		"physiotherapist.serviceAndSchedule.serviceFees": 1,
		"physiotherapist.serviceAndSchedule.slots": bson.M{
			"startTime": 1,
			"endTime":   1,
			"days":      1,
		},
	}

	findOptions := options.FindOne().SetProjection(projection)

	var physiotherapistData entity.ServiceEntity
	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&physiotherapistData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(physiotherapist.GetPhysiotherapistResDto{
			Status:  false,
			Message: "Failed to fetch physiotherapist data: " + err.Error(),
		})
	}

	if physiotherapistData.Physiotherapist == nil {
		return c.Status(fiber.StatusNotFound).JSON(physiotherapist.GetPhysiotherapistResDto{
			Status:  false,
			Message: "physiotherapist data not found",
		})
	}

	scheduleData := make([]physiotherapist.ServiceAndSchedule, 0)
	if physiotherapistData.Physiotherapist != nil && len(physiotherapistData.Physiotherapist.ServiceAndSchedule) > 0 {
		for _, service := range physiotherapistData.Physiotherapist.ServiceAndSchedule {
			var slotData []physiotherapist.Slots
			for _, slots := range service.Slots {
				slotData = append(slotData, physiotherapist.Slots{
					StartTime: slots.StartTime,
					EndTime:   slots.EndTime,
					Days:      slots.Days,
				})
			}
			scheduleData = append(scheduleData, physiotherapist.ServiceAndSchedule{
				Id:          service.Id,
				ServiceFees: service.ServiceFees,
				Name:        service.Name,
				Slots:       slotData,
			})
		}

	}

	physiotherapistRes := physiotherapist.GetPhysiotherapistResDto{
		Status:  true,
		Message: "Physiotherapist data fetched successfully",
		Data: physiotherapist.PhysiotherapistResponse{
			Id:                 physiotherapistData.Id,
			Image:              physiotherapistData.Physiotherapist.Information.Image,
			Name:               physiotherapistData.Physiotherapist.Information.Name,
			AboutMe:            physiotherapistData.Physiotherapist.Information.AdditionalText,
			ServiceAndSchedule: scheduleData,
			TotalReviews:       physiotherapistData.Physiotherapist.Review.TotalReviews,
			AvgRating:          physiotherapistData.Physiotherapist.Review.AvgRating,
		},
	}

	return c.Status(fiber.StatusOK).JSON(physiotherapistRes)
}
