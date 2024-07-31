package customerAuth

import (
	"careville_backend/database"
	"careville_backend/dto/customer/customerAuth"
	"careville_backend/entity"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary Signup provider
// @Description Signup provider
// @Tags customer authorization
// @Param signup body customerAuth.CustomerSignupReqDto true "send 6 digit otp to email for signup"
// @Produce json
// @Success 200 {object} customerAuth.CustomerResponseDto
// @Router /customer/signup [post]
func SignupCustomer(c *fiber.Ctx) error {
	var (
		customerColl = database.GetCollection("customer")
		otpColl      = database.GetCollection("otp")
		data         customerAuth.CustomerSignupReqDto
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(customerAuth.CustomerResponseDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Check if email is not already used
	filter := bson.M{
		"email": strings.ToLower(data.Email),
	}

	exists, err := customerColl.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(500).JSON(customerAuth.CustomerResponseDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	if exists > 0 {
		return c.Status(400).JSON(customerAuth.CustomerResponseDto{
			Status:  false,
			Message: "Email is already in use",
		})
	}

	// Check if dial code and phone number combination is already in use
	filter = bson.M{
		"phoneNumber.dialCode": data.DialCode,
		"phoneNumber.number":   data.PhoneNumber,
	}

	exists, err = customerColl.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(500).JSON(customerAuth.CustomerResponseDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	if exists > 0 {
		return c.Status(400).JSON(customerAuth.CustomerResponseDto{
			Status:  false,
			Message: "phone number with dial code is already in use",
		})
	}

	smallEmail := strings.ToLower(data.Email)

	// otp := utils.Generate6DigitOtp()
	otp := "111111"

	otpData := entity.OtpEntity{
		Id:        primitive.NewObjectID(),
		Otp:       otp,
		Email:     smallEmail,
		CreatedAt: time.Now().UTC(),
	}

	_, err = otpColl.InsertOne(ctx, otpData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(customerAuth.CustomerResponseDto{
			Status:  false,
			Message: "failed to store OTP in the database: " + err.Error(),
		})
	}

	// _, err = utils.SendEmail(data.Email, otp)
	// if err != nil {
	// 	return c.Status(500).JSON(customerAuth.CustomerResponseDto{
	// 		Status:  false,
	// 		Message: "failed to send OTP email: " + err.Error(),
	// 	})
	// }

	return c.Status(fiber.StatusOK).JSON(customerAuth.CustomerResponseDto{
		Status:  true,
		Message: "OTP sent successfully",
	})
}
