package customerAuth

import (
	"careville_backend/database"
	customerAuth "careville_backend/dto/customer/customerAuth"
	"careville_backend/entity"
	"careville_backend/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Reset customer Password after OTP Verification
// @Description Reset customer password after OTP verification using the new password and confirm password
// @Tags customer authorization
// @Accept application/json
// @Param customer body customerAuth.ResetPasswordAfterOtpReqDto true "Reset customer password after OTP verification"
// @Produce json
// @Success 200 {object} customerAuth.CustomerPasswordResDto
// @Router /customer/reset-password [Put]
func ResetPasswordAfterOtp(c *fiber.Ctx) error {
	var (
		customerColl = database.GetCollection("customer")
		data       customerAuth.ResetPasswordAfterOtpReqDto
		provider    entity.ServiceEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(customerAuth.CustomerPasswordResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	filter := bson.M{
		"email": strings.ToLower(data.Email),
	}

	// Find the user with email address from client
	err = customerColl.FindOne(ctx, filter).Decode(&provider)
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
			Message: "Internal server error, while getting the provider: " + err.Error(),
		})
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(data.NewPassword)
	if err != nil {
		return c.Status(500).JSON(customerAuth.CustomerPasswordResDto{
			Status:  false,
			Message: "Failed to hash the password: " + err.Error(),
		})
	}

	_, err = customerColl.UpdateOne(ctx, bson.M{"_id": provider.Id}, bson.M{"$set": bson.M{"password": hashedPassword}})
	if err != nil {
		return c.Status(500).JSON(customerAuth.CustomerPasswordResDto{
			Status:  false,
			Message: "Failed to update password in the database: " + err.Error(),
		})
	}

	return c.Status(200).JSON(customerAuth.CustomerPasswordResDto{
		Status:  true,
		Message: "Password updated successfully after OTP verification",
	})
}
