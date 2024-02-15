package laboratory

import (
	"careville_backend/database"
	laboratory "careville_backend/dto/customer/laboratories"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get all investigations
// @Description Get all investigations
// @Tags customer laboratory
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param laboratoryId path string true "laboratory ID"
// @Produce json
// @Success 200 {object} laboratory.InvestigationResDto
// @Router /customer/healthFacility/get-investigations [get]
func GetInvestigations(c *fiber.Ctx) error {

	var service entity.ServiceEntity

	idParam := c.Query("laboratoryId")
	laboratoryID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(laboratory.GetLaboratoryResDto{
			Status:  false,
			Message: "Invalid laboratory ID",
		})
	}

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"_id": laboratoryID,
	}

	projection := bson.M{
		"laboratory.investigations.id":    1,
		"laboratory.investigations.name":  1,
		"laboratory.investigations.price": 1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(laboratory.InvestigationResDto{
				Status:  false,
				Message: "Investigation not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(laboratory.InvestigationResDto{
			Status:  false,
			Message: "Failed to fetch investigation from MongoDB: " + err.Error(),
		})
	}

	investigationData := make([]laboratory.InvestigationRes, 0)
	if service.Laboratory != nil && len(service.Laboratory.Investigations) > 0 {
		for _, investigation := range service.Laboratory.Investigations {
			investigationData = append(investigationData, laboratory.InvestigationRes{
				Id:    investigation.Id,
				Name:  investigation.Name,
				Price: investigation.Price,
			})
		}
	}

	if len(investigationData) == 0 {
		return c.Status(fiber.StatusOK).JSON(laboratory.InvestigationResDto{
			Status:  false,
			Message: "No Investigation data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(laboratory.InvestigationResDto{
		Status:  true,
		Message: "Investigations retrieved successfully",
		Data:    investigationData,
	})
}
