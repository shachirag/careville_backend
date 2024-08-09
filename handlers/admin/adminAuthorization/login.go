package adminAuth

import (
	"careville_backend/database"
	"careville_backend/dto/admin/adminAuth"
	"careville_backend/entity"
	"careville_backend/utils"
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var ctx = context.Background()

// @Summary send 6 digit otp for Login Admin
// @Description send 6 digit otp for Login Admin
// @Tags admin authorization
// @Accept application/json
// @Param admin body adminAuth.LoginAdminReqDto true "send 6 digit otp for Login Admin"
// @Produce json
// @Success 200 {object} adminAuth.LoginAdminResDto
// @Router /admin/login [post]
func LoginAdmin(c *fiber.Ctx) error {
	var (
		adminColl = database.GetCollection("admin")
		otpColl   = database.GetCollection("otp")
		data      adminAuth.LoginAdminReqDto
		admin     entity.AdminEntity
	)

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(adminAuth.LoginAdminResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	email := strings.ToLower(data.Email)
	err = adminColl.FindOne(ctx, bson.M{"email": email}).Decode(&admin)
	if err != nil {

		if err == mongo.ErrNoDocuments {
			return c.Status(404).JSON(adminAuth.LoginAdminResDto{
				Status:  false,
				Message: "Email does not exists",
			})
		}

		return c.Status(500).JSON(adminAuth.LoginAdminResDto{
			Status:  false,
			Message: "Internal server error while getting the admin: " + err.Error(),
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(strings.TrimSpace(data.Password)))
	if err != nil {
		return c.Status(400).JSON(adminAuth.LoginAdminResDto{
			Status:  false,
			Message: "Invalid credentials",
		})
	}

	// otp := utils.Generate6DigitOtp()
	otp := "111111"
	otpData := entity.OtpEntity{
		Id:        primitive.NewObjectID(),
		Otp:       otp,
		Email:     email,
		CreatedAt: time.Now().UTC(),
	}

	_, err = otpColl.InsertOne(ctx, otpData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(adminAuth.LoginAdminResDto{
			Status:  false,
			Message: "Failed to store OTP in the database: " + err.Error(),
		})
	}

	_, err = utils.SendEmailForPassword(email, otp)
	if err != nil {
		return c.Status(500).JSON(adminAuth.LoginAdminResDto{
			Status:  false,
			Message: "Internal server error, while sending email: " + err.Error(),
		})
	}

	return c.Status(200).JSON(adminAuth.LoginAdminResDto{
		Status:  true,
		Message: "Successfully sent 6-digit OTP.",
	})
}
