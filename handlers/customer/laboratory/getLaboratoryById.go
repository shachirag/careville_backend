package laboratory

import (
	"careville_backend/database"
	laboratory "careville_backend/dto/customer/laboratories"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get laboratory by ID
// @Tags customer laboratory
// @Description Get laboratory by ID
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "laboratory ID"
// @Produce json
// @Success 200 {object} laboratory.GetLaboratoryResDto
// @Router /customer/healthFacility/get-laboratory/{id} [get]
func GetLaboratoryByID(c *fiber.Ctx) error {
	var (
		serviceColl = database.GetCollection("service")
	)

	idParam := c.Params("id")
	laboratoryID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(laboratory.GetLaboratoryResDto{
			Status:  false,
			Message: "Invalid laboratory ID",
		})
	}

	filter := bson.M{"_id": laboratoryID}

	projection := bson.M{
		"laboratory.information.name":           1,
		"laboratory.information.image":          1,
		"_id":                                   1,
		"laboratory.information.additionalText": 1,
		"laboratory.review.totalReviews":        1,
		"laboratory.review.avgRating":           1,
		"laboratory.investigations.id":          1,
		"laboratory.investigations.name":        1,
		"laboratory.investigations.type":        1,
		"laboratory.investigations.information": 1,
		"laboratory.investigations.price":       1,
		"laboratory.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
	}

	findOptions := options.FindOne().SetProjection(projection)

	var laboratoryData entity.ServiceEntity
	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&laboratoryData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(laboratory.GetLaboratoryResDto{
			Status:  false,
			Message: "Failed to fetch laboratory data: " + err.Error(),
		})
	}

	if laboratoryData.Laboratory == nil {
		return c.Status(fiber.StatusNotFound).JSON(laboratory.GetLaboratoryResDto{
			Status:  false,
			Message: "laboratory data not found",
		})
	}

	investigationData := make([]laboratory.Investigations, 0)
	if laboratoryData.Laboratory != nil && len(laboratoryData.Laboratory.Investigations) > 0 {
		for _, investigation := range laboratoryData.Laboratory.Investigations {
			investigationData = append(investigationData, laboratory.Investigations{
				Id:          investigation.Id,
				Type:        investigation.Type,
				Name:        investigation.Name,
				Information: investigation.Information,
				Price:       investigation.Price,
			})
		}
	}

	laboratoryRes := laboratory.GetLaboratoryResDto{
		Status:  true,
		Message: "Laboratory data fetched successfully",
		Data: laboratory.LaboratoryResponse{
			Id:             laboratoryData.Id,
			Image:          laboratoryData.Laboratory.Information.Image,
			Name:           laboratoryData.Laboratory.Information.Name,
			AboutUs:        laboratoryData.Laboratory.Information.AdditionalText,
			Address:        laboratory.Address(laboratoryData.Laboratory.Information.Address),
			Investigations: investigationData,
			TotalReviews:   laboratoryData.Laboratory.Review.TotalReviews,
			AvgRating:      laboratoryData.Laboratory.Review.AvgRating,
		},
	}

	return c.Status(fiber.StatusOK).JSON(laboratoryRes)
}
