package pharmacy

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

// @Summary Add other service
// @Tags pharmacy
// @Description Add other service
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param  provider body services.AddPharmacyOtherServiceReqDto true "add other services"
// @Produce json
// @Success 200 {object} services.UpdatePharmacyOtherServiceResDto
// @Router /provider/services/add-pharmacy-other-service [post]
func AddOtherServices(c *fiber.Ctx) error {
	var (
		servicesColl = database.GetCollection("service")
		data         services.AddPharmacyOtherServiceReqDto
		provider     entity.ServiceEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(services.UpdatePharmacyOtherServiceResDto{
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
				Message: "Provider not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateDoctorImageResDto{
			Status:  false,
			Message: "Failed to fetch provider from MongoDB: " + err.Error(),
		})
	}

	update := bson.M{
		"$addToSet": bson.M{
			"pharmacy.additionalServices": bson.M{
				"$each": []entity.AdditionalServices{
					{
						Id:          primitive.NewObjectID(),
						Name:        data.Name,
						Information: data.Information,
					},
				},
			},
		},
	}

	updateRes, err := servicesColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdatePharmacyOtherServiceResDto{
			Status:  false,
			Message: "Failed to update provider data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(services.UpdatePharmacyOtherServiceResDto{
			Status:  false,
			Message: "Provider not found",
		})
	}

	hospClinicRes := services.UpdatePharmacyOtherServiceResDto{
		Status:  true,
		Message: "Other Service added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(hospClinicRes)
}
