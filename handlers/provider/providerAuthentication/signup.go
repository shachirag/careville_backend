package providerAuthenticate

import (
	"careville_backend/database"
	providerAuth "careville_backend/dto/provider/providerAuth"
	"careville_backend/entity"
	"careville_backend/utils"
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ctx = context.Background()

// @Summary Signup provider
// @Description Signup provider
// @Tags provider authorization
// @Param signup body providerAuth.ProviderSignupReqDto true "send 6 digit otp to email for signup"
// @Produce json
// @Success 200 {object} providerAuth.ProviderResponseDto
// @Router /provider/signup [post]
func SignupProvider(c *fiber.Ctx) error {
	var (
		serviceColl = database.GetCollection("service")
		otpColl     = database.GetCollection("otp")
		data        providerAuth.ProviderSignupReqDto
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(providerAuth.ProviderResponseDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Check if email is not already used
	filter := bson.M{
		"email": strings.ToLower(data.Email),
	}

	exists, err := serviceColl.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(500).JSON(providerAuth.ProviderResponseDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	if exists > 0 {
		return c.Status(400).JSON(providerAuth.ProviderResponseDto{
			Status:  false,
			Message: "Email is already in use",
		})
	}

	// Check if dial code and phone number combination is already in use
	filter = bson.M{
		"phoneNumber.dialCode": data.DialCode,
		"phoneNumber.number":   data.PhoneNumber,
	}

	exists, err = serviceColl.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(500).JSON(providerAuth.ProviderResponseDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	if exists > 0 {
		return c.Status(400).JSON(providerAuth.ProviderResponseDto{
			Status:  false,
			Message: "phone number with dial code is already in use",
		})
	}

	smallEmail := strings.ToLower(data.Email)

	// Generate a 6-digit OTP
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
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.ProviderResponseDto{
			Status:  false,
			Message: "failed to store OTP in the database: " + err.Error(),
		})
	}
	// Send OTP to the provided email
	_, err = utils.SendEmail(data.Email, otp)
	if err != nil {
		return c.Status(500).JSON(providerAuth.ProviderResponseDto{
			Status:  false,
			Message: "failed to send OTP email: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(providerAuth.ProviderResponseDto{
		Status:  true,
		Message: "OTP sent successfully",
	})
}
