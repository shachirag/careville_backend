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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get other service info
// @Description Get other service info
// @Tags fitnessCenter
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param subscriptionId path string true "other subscription ID"
// @Produce json
// @Success 200 {object} services.GetSubscriptionResDto
// @Router /provider/services/get-subscription-info/{subscriptionId} [get]
func GetSubscriptionInfo(c *fiber.Ctx) error {

	var service entity.ServiceEntity

	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	subscriptionId := c.Params("subscriptionId")
	subscriptionObjID, err := primitive.ObjectIDFromHex(subscriptionId)

	if err != nil {
		return c.Status(400).JSON(services.GetSubscriptionResDto{
			Status:  false,
			Message: "invalid objectId " + err.Error(),
		})
	}

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"_id": providerData.ProviderId,
		"fitnessCenter.subscription": bson.M{
			"$elemMatch": bson.M{
				"id": subscriptionObjID,
			},
		},
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
			return c.Status(fiber.StatusNotFound).JSON(services.GetSubscriptionResDto{
				Status:  false,
				Message: "service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.GetSubscriptionResDto{
			Status:  false,
			Message: "Failed to fetch service from MongoDB: " + err.Error(),
		})
	}

	if service.FitnessCenter == nil {
		return c.Status(fiber.StatusNotFound).JSON(services.GetSubscriptionResDto{
			Status:  false,
			Message: "No subscription information found for the service",
		})
	}

	var subscriptionsRes services.SubscriptionRes

	for _, services1 := range service.FitnessCenter.Subscription {
		if services1.Id == subscriptionObjID {
			subscriptionRes := services.SubscriptionRes{
				Id:      services1.Id,
				Type:    services1.Type,
				Price:   services1.Price,
				Details: services1.Details,
			}

			subscriptionsRes = subscriptionRes
			break
		}
	}

	return c.Status(fiber.StatusOK).JSON(services.GetSubscriptionResDto{
		Status:  true,
		Message: "other service retrieved successfully",
		Data:    subscriptionsRes,
	})
}
