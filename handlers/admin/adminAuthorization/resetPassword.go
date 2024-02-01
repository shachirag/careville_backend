package adminAuth

import (
	"careville_backend/database"
	adminAuth "careville_backend/dto/admin/adminAuth"
	"careville_backend/entity"
	"careville_backend/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Reset admin Password after OTP Verification
// @Description Reset admin password after OTP verification using the new password and confirm password
// @Tags admin authorization
// @Accept application/json
// @Param admin body adminAuth.ResetPasswordAfterOtpReqDto true "Reset admin password after OTP verification"
// @Produce json
// @Success 200 {object} adminAuth.AdminPasswordResDto
// @Router /admin/reset-password [Put]
func ResetPasswordAfterOtp(c *fiber.Ctx) error {
	var (
		adminColl = database.GetCollection("admin")
		data      adminAuth.ResetPasswordAfterOtpReqDto
		provider  entity.ServiceEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(adminAuth.AdminPasswordResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	filter := bson.M{
		"email": strings.ToLower(data.Email),
	}

	// Find the user with email address from client
	err = adminColl.FindOne(ctx, filter).Decode(&provider)
	if err != nil {
		// Check if there is no documents found error
		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(adminAuth.AdminPasswordResDto{
				Status:  false,
				Message: "No admin found",
			})
		}

		return c.Status(500).JSON(adminAuth.AdminPasswordResDto{
			Status:  false,
			Message: "Internal server error, while getting the admin: " + err.Error(),
		})
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(data.NewPassword)
	if err != nil {
		return c.Status(500).JSON(adminAuth.AdminPasswordResDto{
			Status:  false,
			Message: "Failed to hash the password: " + err.Error(),
		})
	}

	_, err = adminColl.UpdateOne(ctx, bson.M{"_id": provider.Id}, bson.M{"$set": bson.M{"password": hashedPassword}})
	if err != nil {
		return c.Status(500).JSON(adminAuth.AdminPasswordResDto{
			Status:  false,
			Message: "Failed to update password in the database: " + err.Error(),
		})
	}

	return c.Status(200).JSON(adminAuth.AdminPasswordResDto{
		Status:  true,
		Message: "Password updated successfully after OTP verification",
	})
}
