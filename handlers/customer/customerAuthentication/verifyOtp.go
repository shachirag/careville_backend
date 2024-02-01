package customerAuth

import (
	"careville_backend/database"
	"careville_backend/dto/customer/customerAuth"
	"careville_backend/entity"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Verify customer OTP
// @Description Verify the entered 6 digit OTP
// @Tags customer authorization
// @Accept application/json
// @Param customer body customerAuth.VerifyOtpReqDto true "Verify 6 digit OTP"
// @Produce json
// @Success 200 {object} customerAuth.CustomerPasswordResDto
// @Router /customer/verify-otp [post]
func VerifyOtp(c *fiber.Ctx) error {
	var (
		otpColl = database.GetCollection("otp")
		data    customerAuth.VerifyOtpReqDto
		otpData entity.OtpEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(customerAuth.CustomerPasswordResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Error handling
	if data.EnteredOTP == "" {
		return c.Status(400).JSON(customerAuth.CustomerPasswordResDto{
			Status:  false,
			Message: "OTP is required",
		})
	}

	smallEmail := strings.ToLower(data.Email)

	// Find the user with email address from client
	err = otpColl.FindOne(ctx, bson.M{"email": smallEmail}, options.FindOne().SetSort(bson.M{"createdAt": -1})).Decode(&otpData)
	if err != nil {
		// Check if there is no documents found error
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(customerAuth.CustomerPasswordResDto{
				Status:  false,
				Message: "Invalid OTP",
			})
		}

		return c.Status(500).JSON(customerAuth.CustomerPasswordResDto{
			Status:  false,
			Message: "Internal server error, while getting the provider: " + err.Error(),
		})
	}

	if data.EnteredOTP != otpData.Otp {

		return c.Status(400).JSON(customerAuth.CustomerPasswordResDto{
			Status:  false,
			Message: "Invalid OTP",
		})
	}

	return c.Status(200).JSON(customerAuth.CustomerPasswordResDto{
		Status:  true,
		Message: "OTP verified successfully",
	})
}
