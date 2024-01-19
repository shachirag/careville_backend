package providerAuthenticate

import (
	"careville_backend/database"
	providerAuth "careville_backend/dto/provider/providerAuth"
	"careville_backend/entity"
	"careville_backend/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Forgot Password
// @Description Forgot Password
// @Tags provider authorization
// @Accept application/json
// @Param provider body providerAuth.ProviderForgotPasswordReqDto true "forgot password for provider"
// @Produce json
// @Success 200 {object} providerAuth.ProviderPasswordResDto
// @Router /provider/forgot-password [post]
func ForgotPassword(c *fiber.Ctx) error {
	var (
		serviceColl = database.GetCollection("service")
		otpColl     = database.GetCollection("otp")
		data        providerAuth.ProviderForgotPasswordReqDto
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

	smallEmail := strings.ToLower(data.Email)

	// Find the user with email address from client
	err = serviceColl.FindOne(ctx, bson.M{"user.email": smallEmail}).Decode(&provider)
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

	// Generate 6-digit OTP
	otp := utils.Generate6DigitOtp()

	// Store the OTP in the forgotPassword collection
	otpData := entity.OtpEntity{
		Id:        primitive.NewObjectID(),
		Otp:       otp,
		Email:     smallEmail,
		CreatedAt: time.Now().UTC(),
	}

	_, err = otpColl.InsertOne(ctx, otpData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.ProviderPasswordResDto{
			Status:  false,
			Message: "Failed to store OTP in the database: " + err.Error(),
		})
	}

	// Sending email to the recipient with the OTP
	_, err = utils.SendEmailForPassword(data.Email, otp)
	if err != nil {
		return c.Status(500).JSON(providerAuth.ProviderPasswordResDto{
			Status:  false,
			Message: "Internal server error, while sending email: " + err.Error(),
		})
	}

	return c.Status(200).JSON(providerAuth.ProviderPasswordResDto{
		Status:  true,
		Message: "Successfully sent 6 digit OTP.",
	})
}
