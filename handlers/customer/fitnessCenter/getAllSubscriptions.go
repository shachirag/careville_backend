package fitnessCenter

import (
	"careville_backend/database"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get all subscriptions
// @Description Get all subscriptions
// @Tags customer fitnessCenter
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "fitnessCenter ID"
//
// @Produce json
// @Success 200 {object} services.SubscriptionResDto
// @Router /customer/healthFacility/get-subscriptions/{id} [get]
func GetAllSubscriptions(c *fiber.Ctx) error {

	var service entity.ServiceEntity

	serviceColl := database.GetCollection("service")

	idParam := c.Params("id")
	fitnessCenterID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.SubscriptionResDto{
			Status:  false,
			Message: "Invalid fitnessCenter ID",
		})
	}

	filter := bson.M{"_id": fitnessCenterID}

	filter = bson.M{
		"_id": fitnessCenterID,
	}

	projection := bson.M{
		"fitnessCenter.subscription.id":      1,
		"fitnessCenter.subscription.type":    1,
		"fitnessCenter.subscription.details": 1,
		"fitnessCenter.subscription.price":   1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
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
