package nurse

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary provider currently available
// @Tags nurse
// @Description currently available
// @Accept json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param isEmergencyAvailable query bool false "isEmergencyAvailable value (true or false)"
// @Produce json
// @Success 200 {object} services.NotificationResDto
// @Router /provider/services/nurse-currently-available [put]
func ProviderNurseCurrentlyAvailable(c *fiber.Ctx) error {

	providerData := providerMiddleware.GetProviderMiddlewareData(c)
	isEmergencyAvailableParam := c.Query("isEmergencyAvailable")
	isEnable, err := strconv.ParseBool(isEmergencyAvailableParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.NotificationResDto{
			Status:  false,
			Message: "Invalid value for isEnable: " + err.Error(),
		})
	}

	serviceColl := database.GetCollection("service")

	filter := bson.M{"_id": providerData.ProviderId}

	update := bson.M{"$set": bson.M{
		"physiotherapist.information.isEmergencyAvailable": isEnable,
	},
	}

	opts := options.Update().SetUpsert(true)
	result, err := serviceColl.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.NotificationResDto{
			Status:  false,
			Message: "Failed to change currently avalilable: " + err.Error(),
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
