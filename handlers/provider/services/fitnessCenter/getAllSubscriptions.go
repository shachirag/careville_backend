package fitnessCenter

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get all subscriptions
// @Description Get all subscriptions
// @Tags fitnessCenter
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} services.SubscriptionResDto
// @Router /provider/services/get-all-subscriptions [get]
func GetAllSubscriptions(c *fiber.Ctx) error {
	ctx := context.Background()

	var service entity.ServiceEntity

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"_id": providerData.ProviderId,
	}

	projection := bson.M{
		"fitnessCenter.subscription.id":      1,
		"fitnessCenter.subscription.type":    1,
		"fitnessCenter.subscription.details": 1,
		"fitnessCenter.subscription.price":   1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	err := serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.SubscriptionResDto{
				Status:  false,
				Message: "service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.SubscriptionResDto{
			Status:  false,
			Message: "Failed to fetch service from MongoDB: " + err.Error(),
		})
	}

	subscriptionData := make([]services.SubscriptionRes, 0)
	if service.FitnessCenter != nil && len(service.FitnessCenter.Subscription) > 0 {
		for _, subscription := range service.FitnessCenter.Subscription {
			subscriptionData = append(subscriptionData, services.SubscriptionRes{
				Id:      subscription.Id,
				Type:    subscription.Type,
				Details: subscription.Details,
				Price:   subscription.Price,
			})
		}
	}

	if len(subscriptionData) == 0 {
		return c.Status(fiber.StatusOK).JSON(services.SubscriptionResDto{
			Status:  false,
			Message: "No subscription data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(services.SubscriptionResDto{
		Status:  true,
		Message: "subscriptions retrieved successfully",
		Data:    subscriptionData,
	})
}
