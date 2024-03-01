package adminAuth

import (
	"careville_backend/database"
	adminAuth "careville_backend/dto/admin/adminAuth"
	"careville_backend/entity"
	"careville_backend/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Resend OTP
// @Description Resend 6 digit OTP to email
// @Tags admin authorization
// @Accept application/json
// @Param user body adminAuth.ResendOtpReqDto true "Resend 6 digit OTP to email"
// @Produce json
// @Success 200 {object} adminAuth.UserPasswordResDto
// @Router /admin/resend-otp [post]
func ResendOTP(c *fiber.Ctx) error {
	var (
		otpColl = database.GetCollection("otp")
		data    adminAuth.ResendOtpReqDto
		otpData entity.OtpEntity
	)

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(adminAuth.UserPasswordResDto{
			Status:  false,
			Message: err.Error(),
		})
	}
	smallEmail := strings.ToLower(data.Email)

	err = otpColl.FindOne(ctx, bson.M{"email": smallEmail}).Decode(&otpData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(adminAuth.UserPasswordResDto{
				Status:  false,
				Message: "No OTP found for the provided email. Please request a new OTP.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(adminAuth.UserPasswordResDto{
			Status:  false,
			Message: "Internal server error while fetching OTP data: " + err.Error(),
		})
	}

	// Generate a new 6-digit OTP
	newOTP := utils.Generate6DigitOtp()

	// Update the existing OTP with the new OTP and reset the creation time
	otpData.Otp = newOTP
	otpData.CreatedAt = time.Now().UTC()

	update := bson.M{"$set": bson.M{"otp": newOTP, "createdAt": otpData.CreatedAt}}

	_, err = otpColl.UpdateOne(ctx, bson.M{"email": smallEmail}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(adminAuth.UserPasswordResDto{
			Status:  false,
			Message: "Failed to update OTP in the database: " + err.Error(),
		})
	}

	_, err = utils.SendEmail(data.Email, newOTP)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(adminAuth.UserPasswordResDto{
			Status:  false,
			Message: "Internal server error while resending email: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(adminAuth.UserPasswordResDto{
		Status:  true,
		Message: "Successfully resent 6 digit OTP.",
	})
}
