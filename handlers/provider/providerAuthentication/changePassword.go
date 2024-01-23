package providerAuthenticate

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	providerAuth "careville_backend/dto/provider/providerAuth"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// ChangeProviderPassword is the handler for changing provider passwords
// @Summary Change provider Password
// @Description Change provider Password
// @Tags provider authorization
// @Accept application/json
// @Param Authorization header string true "Authentication header"
// @Param provider body providerAuth.ProviderChangePasswordReqDto true "Change password of provider"
// @Produce json
// @Success 200 {object} providerAuth.ProviderChangePasswordResDto
// @Router /provider/profile/change-password [put]
func ChangeProviderPassword(c *fiber.Ctx) error {
	var (
		serviceColl = database.GetCollection("service")
		data        providerAuth.ProviderChangePasswordReqDto
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.ProviderChangePasswordResDto{
			Status:  false,
			Message: "Failed to parse request body: " + err.Error(),
		})
	}

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	filter := bson.M{"_id": providerData.ProviderId}

	result := serviceColl.FindOne(c.Context(), filter)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(providerAuth.ProviderChangePasswordResDto{
				Status:  false,
				Message: "Provider not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.ProviderChangePasswordResDto{
			Status:  false,
			Message: "Error finding provider: " + err.Error(),
		})
	}

	// Decode the provider data
	var provider entity.ServiceEntity
	err = result.Decode(&provider)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.ProviderChangePasswordResDto{
			Status:  false,
			Message: "Failed to decode provider data: " + err.Error(),
		})
	}

	if data.CurrentPassword != "" {
		// Validate the current password
		err = bcrypt.CompareHashAndPassword([]byte(provider.User.Password), []byte(data.CurrentPassword))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(providerAuth.ProviderChangePasswordResDto{
				Status:  false,
				Message: "Current password is incorrect",
			})
		}
	}

	// Hash the new password
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), 6)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.ProviderChangePasswordResDto{
			Status:  false,
			Message: "Failed to hash the new password",
		})
	}

	update := bson.M{
		"$set": bson.M{
			"user.password": string(hashedNewPassword),
		},
	}

	// Execute the update operation
	updateRes, err := serviceColl.UpdateOne(c.Context(), filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.ProviderChangePasswordResDto{
			Status:  false,
			Message: "Failed to update provider password in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(providerAuth.ProviderChangePasswordResDto{
			Status:  false,
			Message: "Provider not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(providerAuth.ProviderChangePasswordResDto{
		Status:  true,
		Message: "Provider password updated successfully",
	})
}
