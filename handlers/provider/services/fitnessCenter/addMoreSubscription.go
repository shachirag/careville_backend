package fitnessCenter

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

// @Summary Add more subscription
// @Tags fitnessCenter
// @Description Add more subscription
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param  provider body services.SubscriptionReqDto true "add trainer"
// @Produce json
// @Success 200 {object} services.TrainerResponseDto
// @Router /provider/services/add-more-subscription [post]
func AddMoreSubscriptions(c *fiber.Ctx) error {
	var (
		servicesColl = database.GetCollection("service")
		data         services.SubscriptionReqDto
		provider     entity.ServiceEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(services.SubscriptionResponseDto{
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
			return c.Status(fiber.StatusNotFound).JSON(services.SubscriptionResponseDto{
				Status:  false,
				Message: "provider not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.SubscriptionResponseDto{
			Status:  false,
			Message: "Failed to fetch provider from MongoDB: " + err.Error(),
		})
	}

	update := bson.M{
		"$addToSet": bson.M{
			"fitnessCenter.subscription": bson.M{
				"$each": []entity.Subscription{
					{
						Id:      primitive.NewObjectID(),
						Type:    data.Type,
						Details: data.Details,
						Price:   data.Price,
					},
				},
			},
		},
	}

	updateRes, err := servicesColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.TrainerResponseDto{
			Status:  false,
			Message: "Failed to update provider data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(services.TrainerResponseDto{
			Status:  false,
			Message: "provider not found",
		})
	}

	hospClinicRes := services.TrainerResponseDto{
		Status:  true,
		Message: "Subscription added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(hospClinicRes)
}
