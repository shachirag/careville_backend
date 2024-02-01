package customerAuth

import (
	"careville_backend/database"
	customerAuth "careville_backend/dto/customer/customerAuth"
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
// @Tags customer authorization
// @Accept application/json
// @Param customer body customerAuth.CustomerForgotPasswordReqDto true "forgot password for customer"
// @Produce json
// @Success 200 {object} customerAuth.CustomerPasswordResDto
// @Router /customer/forgot-password [post]
func ForgotPassword(c *fiber.Ctx) error {
	var (
		customerColl = database.GetCollection("customer")
		otpColl     = database.GetCollection("otp")
		data        customerAuth.CustomerForgotPasswordReqDto
		customer    entity.CustomerEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(customerAuth.CustomerPasswordResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	smallEmail := strings.ToLower(data.Email)

	// Find the user with email address from client
	err = customerColl.FindOne(ctx, bson.M{"email": smallEmail}).Decode(&customer)
	if err != nil {
		// Check if there is no documents found error
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(customerAuth.CustomerPasswordResDto{
				Status:  false,
				Message: "No customer found",
			})
		}

		return c.Status(500).JSON(customerAuth.CustomerPasswordResDto{
			Status:  false,
			Message: "Internal server error, while getting the customer: " + err.Error(),
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
		return c.Status(fiber.StatusInternalServerError).JSON(customerAuth.CustomerPasswordResDto{
			Status:  false,
			Message: "Failed to store OTP in the database: " + err.Error(),
		})
	}

	_, err = utils.SendEmailForPassword(data.Email, otp)
	if err != nil {
		return c.Status(500).JSON(customerAuth.CustomerPasswordResDto{
			Status:  false,
			Message: "Internal server error, while sending email: " + err.Error(),
		})
	}

	return c.Status(200).JSON(customerAuth.CustomerPasswordResDto{
		Status:  true,
		Message: "Successfully sent 6 digit OTP.",
	})
}
