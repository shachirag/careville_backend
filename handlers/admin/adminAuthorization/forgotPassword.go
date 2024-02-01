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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Forgot Password
// @Description Forgot Password
// @Tags admin authorization
// @Accept application/json
// @Param admin body adminAuth.LoginAdminReqDto true "forgot password for customer"
// @Produce json
// @Success 200 {object} adminAuth.AdminPasswordResDto
// @Router /admin/forgot-password [post]
func ForgotPassword(c *fiber.Ctx) error {
	var (
		adminColl = database.GetCollection("admin")
		otpColl   = database.GetCollection("otp")
		data      adminAuth.LoginAdminReqDto
		admin     entity.AdminEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(adminAuth.AdminPasswordResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	smallEmail := strings.ToLower(data.Email)

	// Find the user with email address from client
	err = adminColl.FindOne(ctx, bson.M{"email": smallEmail}).Decode(&admin)
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
		return c.Status(fiber.StatusInternalServerError).JSON(adminAuth.AdminPasswordResDto{
			Status:  false,
			Message: "Failed to store OTP in the database: " + err.Error(),
		})
	}

	_, err = utils.SendEmailForPassword(data.Email, otp)
	if err != nil {
		return c.Status(500).JSON(adminAuth.AdminPasswordResDto{
			Status:  false,
			Message: "Internal server error, while sending email: " + err.Error(),
		})
	}

	return c.Status(200).JSON(adminAuth.AdminPasswordResDto{
		Status:  true,
		Message: "Successfully sent 6 digit OTP.",
	})
}
