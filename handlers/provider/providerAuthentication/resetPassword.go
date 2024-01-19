package providerAuthenticate

import (
	"careville_backend/database"
	providerAuth "careville_backend/dto/provider/providerAuth"
	"careville_backend/entity"
	"careville_backend/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Reset provider Password after OTP Verification
// @Description Reset provider password after OTP verification using the new password and confirm password
// @Tags provider authorization
// @Accept application/json
// @Param provider body providerAuth.ResetPasswordAfterOtpReqDto true "Reset provider password after OTP verification"
// @Produce json
// @Success 200 {object} providerAuth.ProviderPasswordResDto
// @Router /provider/reset-password [Put]
func ResetPasswordAfterOtp(c *fiber.Ctx) error {
	var (
		serviceColl = database.GetCollection("service")
		data        providerAuth.ResetPasswordAfterOtpReqDto
		provider    entity.ServiceEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(providerAuth.ProviderPasswordResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	filter := bson.M{
		"user.email": strings.ToLower(data.Email),
	}

	// Find the user with email address from client
	err = serviceColl.FindOne(ctx, filter).Decode(&provider)
	if err != nil {
		// Check if there is no documents found error
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(providerAuth.ProviderPasswordResDto{
				Status:  false,
				Message: "No provider found",
			})
		}

		return c.Status(500).JSON(providerAuth.ProviderPasswordResDto{
			Status:  false,
			Message: "Internal server error, while getting the provider: " + err.Error(),
		})
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(data.NewPassword)
	if err != nil {
		return c.Status(500).JSON(providerAuth.ProviderPasswordResDto{
			Status:  false,
			Message: "Failed to hash the password: " + err.Error(),
		})
	}

	_, err = serviceColl.UpdateOne(ctx, bson.M{"_id": provider.Id}, bson.M{"$set": bson.M{"user.password": hashedPassword}})
	if err != nil {
		return c.Status(500).JSON(providerAuth.ProviderPasswordResDto{
			Status:  false,
			Message: "Failed to update password in the database: " + err.Error(),
		})
	}

	return c.Status(200).JSON(providerAuth.ProviderPasswordResDto{
		Status:  true,
		Message: "Password updated successfully after OTP verification",
	})
}
