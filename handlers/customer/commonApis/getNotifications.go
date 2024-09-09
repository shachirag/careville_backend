package common

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get Notifications
// @Tags customer commonApis
// @Description Get Notifications
// @Produce json
// @Success 200 {object} services.NotificationResData
// @Router /customer/get-notifications [get]
func GetAllNotifications(c *fiber.Ctx) error {

	var (
		notificationColl = database.GetCollection("notification")
	)

	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	filter := bson.M{
		"customerId": providerData.ProviderId,
	}

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	cursor, err := notificationColl.Find(ctx, filter, sortOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.NotificationResData{
			Status:  false,
			Message: "Failed to fetch notifications data: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var notificationData []services.GetNotificationRes
	for cursor.Next(ctx) {
		var notification entity.NotificationEntity
		if err := cursor.Decode(&notification); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(services.NotificationResData{
				Status:  false,
				Message: "Failed to decode notification data: " + err.Error(),
			})
		}

		notificationData = append(notificationData, services.GetNotificationRes{
			Id:        notification.Id,
			Title:     notification.Title,
			Body:      notification.Body,
			Data:      notification.Data,
			CreatedAt: notification.CreatedAt,
			UpdatedAt: notification.UpdatedAt,
		})
	}

	if len(notificationData) == 0 {
		return c.Status(fiber.StatusOK).JSON(services.NotificationResData{
			Status:  true,
			Message: "No Notification found.",
			Data:    []services.GetNotificationRes{},
		})
	}

	return c.Status(fiber.StatusOK).JSON(services.NotificationResData{
		Status:  true,
		Message: "Successfully fetched Notifications",
		Data:    notificationData,
	})
}
