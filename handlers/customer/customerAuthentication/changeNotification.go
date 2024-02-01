package customerAuth

import (
	"careville_backend/database"
	customerMiddleware "careville_backend/dto/customer/middleware"
	"careville_backend/dto/provider/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// @Summary customer notification
// @Tags customer authorization
// @Description customer notification
// @Accept json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param isEnable query bool false "isEnable value (true or false)"
// @Produce json
// @Success 200 {object} services.NotificationResDto
// @Router /customer/profile/change-notification [put]
func CustomerNotification(c *fiber.Ctx) error {

	customerData := customerMiddleware.GetCustomerMiddlewareData(c)
	isEnableParam := c.Query("isEnable")
	isEnable, err := strconv.ParseBool(isEnableParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.NotificationResDto{
			Status:  false,
			Message: "Invalid value for isEnable: " + err.Error(),
		})
	}

	customerColl := database.GetCollection("customer")

	filter := bson.M{"_id": customerData.CustomerId}

	update := bson.M{"$set": bson.M{"notification.isEnabled": isEnable}}

	result, err := customerColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.NotificationResDto{
			Status:  false,
			Message: "Failed to change notification: " + err.Error(),
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
