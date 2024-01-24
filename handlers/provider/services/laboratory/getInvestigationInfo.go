package laboratory

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

// @Summary Get investigation info
// @Description Get investigation info
// @Tags laboratory
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param investigationId path string true "investigation ID"
// @Produce json
// @Success 200 {object} services.DoctorResDto
// @Router /provider/services/get-investigation-info/{investigationId} [get]
func GetInvesitagtionInfo(c *fiber.Ctx) error {
	ctx := context.TODO()

	var service entity.ServiceEntity

	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	investigationId := c.Params("investigationId")
	investigationObjID, err := primitive.ObjectIDFromHex(investigationId)

	if err != nil {
		return c.Status(400).JSON(services.GetInvestigationResDto{
			Status:  false,
			Message: "invalid objectId " + err.Error(),
		})
	}

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"_id": providerData.ProviderId,
		"laboratory.investigations": bson.M{
			"$elemMatch": bson.M{
				"id": investigationObjID,
			},
		},
	}

	projection := bson.M{
		"laboratory.investigations.id":          1,
		"laboratory.investigations.name":        1,
		"laboratory.investigations.type":        1,
		"laboratory.investigations.information": 1,
		"laboratory.investigations.price":       1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.GetInvestigationResDto{
				Status:  false,
				Message: "service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.GetInvestigationResDto{
			Status:  false,
			Message: "Failed to fetch service from MongoDB: " + err.Error(),
		})
	}

	if service.Laboratory == nil {
		return c.Status(fiber.StatusNotFound).JSON(services.GetInvestigationResDto{
			Status:  false,
			Message: "No Laboratory information found for the service",
		})
	}

	var investigationsRes services.InvestigationRes

	for _, investigation := range service.Laboratory.Investigations {
		if investigation.Id == investigationObjID {
			investigationRes := services.InvestigationRes{
				Id:          investigation.Id,
				Name:        investigation.Name,
				Type:        investigation.Type,
				Information: investigation.Information,
				Price:       investigation.Price,
			}

			investigationsRes = investigationRes
			break
		}
	}

	return c.Status(fiber.StatusOK).JSON(services.GetInvestigationResDto{
		Status:  true,
		Message: "investigation retrieved successfully",
		Data:    investigationsRes,
	})
}
