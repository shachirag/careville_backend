package providerAuthenticate

import (
	"careville_backend/database"
	providerAuth "careville_backend/dto/provider/providerAuth"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Verify provider OTP
// @Description Verify the entered 6 digit OTP
// @Tags provider authorization
// @Accept application/json
// @Param provider body providerAuth.VerifyOtpReqDto true "Verify 6 digit OTP"
// @Produce json
// @Success 200 {object} providerAuth.ProviderPasswordResDto
// @Router /provider/verify-otp [post]
func VerifyOtp(c *fiber.Ctx) error {
	var (
		otpColl = database.GetCollection("otp")
		data    providerAuth.VerifyOtpReqDto
		otpData entity.OtpEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(providerAuth.ProviderPasswordResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Error handling
	if data.EnteredOTP == "" {
		return c.Status(400).JSON(providerAuth.ProviderPasswordResDto{
			Status:  false,
			Message: "Entered OTP is required",
		})
	}

	// Find the user with email address from client
	err = otpColl.FindOne(ctx, bson.M{"email": data.Email}, options.FindOne().SetSort(bson.M{"createdAt": -1})).Decode(&otpData)
	if err != nil {
		// Check if there is no documents found error
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(providerAuth.ProviderPasswordResDto{
				Status:  false,
				Message: "Invalid OTP",
			})
		}

		return c.Status(500).JSON(providerAuth.ProviderPasswordResDto{
			Status:  false,
			Message: "Internal server error, while getting the provider: " + err.Error(),
		})
	}

	// Compare the entered OTP with the OTP from the database
	if data.EnteredOTP != otpData.Otp {

		return c.Status(400).JSON(providerAuth.ProviderPasswordResDto{
			Status:  false,
			Message: "Invalid OTP",
		})
	}

	return c.Status(200).JSON(providerAuth.ProviderPasswordResDto{
		Status:  true,
		Message: "OTP verified successfully",
	})
}
