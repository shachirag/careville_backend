package services

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Add more investigations
// @Tags services
// @Description Add more investigations
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param  provider body services.HospitalClinicRequestDto true "add investigation"
// @Produce json
// @Success 200 {object} services.InvestigationResponseDto
// @Router /provider/services/add-more-investigation [post]
func AddMoreInvestigstions(c *fiber.Ctx) error {
	var (
		servicesColl = database.GetCollection("service")
		data         services.InvestigationReqDto
		provider     entity.ServiceEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(services.InvestigationResponseDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	filter := bson.M{
		"_id": providerData.ProviderId,
	}

	err = servicesColl.FindOne(ctx, filter).Decode(&provider)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.UpdateDoctorImageResDto{
				Status:  false,
				Message: "provider not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateDoctorImageResDto{
			Status:  false,
			Message: "Failed to fetch provider from MongoDB: " + err.Error(),
		})
	}

	update := bson.M{
		"$addToSet": bson.M{
			"laboratory.investigations": bson.M{
				"$each": []entity.Investigations{
					{
						Id:          primitive.NewObjectID(),
						Name:        data.Name,
						Type:        data.Type,
						Information: data.Information,
						Price:       data.Price,
					},
				},
			},
		},
	}

	updateRes, err := servicesColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.InvestigationResponseDto{
			Status:  false,
			Message: "Failed to update provider data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(services.InvestigationResponseDto{
			Status:  false,
			Message: "provider not found",
		})
	}

	hospClinicRes := services.InvestigationResponseDto{
		Status:  true,
		Message: "investigation added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(hospClinicRes)
}
