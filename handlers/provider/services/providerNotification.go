package services

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// @Summary provider notification
// @Tags recovery homes in subadmin pannel
// @Description provider notification
// @Accept json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param isEnable query bool false "isEnable value (true or false)"
// @Produce json
// @Success 200 {object} services.NotificationResDto
// @Router /provider/services/change-notification [put]
func ProviderNotification(c *fiber.Ctx) error {

	providerData := providerMiddleware.GetProviderMiddlewareData(c)
	isDeletedParam := c.Query("isEnable")
	isEnable, err := strconv.ParseBool(isDeletedParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.NotificationResDto{
			Status:  false,
			Message: "Invalid value for isEnable: " + err.Error(),
		})
	}

	serviceColl := database.GetCollection("service")

	filter := bson.M{"_id": providerData.ProviderId}

	update := bson.M{"$set": bson.M{"user.notification.isEnable": isEnable}}

	result, err := serviceColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.NotificationResDto{
			Status:  false,
			Message: "Failed tochange notification: " + err.Error(),
		})
	}

	if result.ModifiedCount == 0 {
		return c.Status(fiber.StatusOK).JSON(services.NotificationResDto{
			Status:  false,
			Message: "No documents were modified",
		})
	}

	response := services.NotificationResDto{
		Status:  true,
		Message: "success",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
