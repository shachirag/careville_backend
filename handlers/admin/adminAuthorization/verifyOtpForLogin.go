package adminAuth

import (
	"careville_backend/database"
	"careville_backend/dto/admin/adminAuth"
	"careville_backend/dto/provider/providerAuth"
	"careville_backend/entity"
	"os"
	"strings"
	"time"

	jtoken "github.com/golang-jwt/jwt/v4"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Verify OTP for login
// @Description Verify the entered 6 digit OTP
// @Tags admin authorization
// @Accept application/json
// @Param admin body adminAuth.LoginVerifyOtpReqDto true "Verify 6 digit OTP and then login"
// @Produce json
// @Success 200 {object} adminAuth.LoginVerifyOtpResDto
// @Router /admin/verify-otp-for-login [post]
func VerifyOtpForLogin(c *fiber.Ctx) error {
	var (
		otpColl   = database.GetCollection("otp")
		adminColl = database.GetCollection("admin")
		data      adminAuth.LoginVerifyOtpReqDto
		otpData   entity.OtpEntity
		admin     entity.AdminEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(adminAuth.LoginVerifyOtpResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Error handling for empty entered OTP
	if data.EnteredOTP == "" {
		return c.Status(400).JSON(adminAuth.LoginVerifyOtpResDto{
			Status:  false,
			Message: "Entered OTP is required",
		})
	}

	email := strings.ToLower(data.Email)
	// Find the user with email address from client in OTP collection
	err = otpColl.FindOne(ctx, bson.M{"email": email}, options.FindOne().SetSort(bson.M{"createdAt": -1})).Decode(&otpData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(adminAuth.LoginVerifyOtpResDto{
				Status:  false,
				Message: "Invalid OTP",
			})
		}
		return c.Status(500).JSON(adminAuth.LoginVerifyOtpResDto{
			Status:  false,
			Message: "Internal server error, while getting the user: " + err.Error(),
		})
	}

	// Compare the entered OTP with the OTP from the database
	if data.EnteredOTP != otpData.Otp {
		return c.Status(400).JSON(adminAuth.LoginVerifyOtpResDto{
			Status:  false,
			Message: "Invalid OTP",
		})
	}

	err = adminColl.FindOne(ctx, bson.M{"email": data.Email}).Decode(&admin)
	if err != nil {
		return c.Status(400).JSON(adminAuth.LoginVerifyOtpResDto{
			Status:  false,
			Message: "Admin not found or internal server error: " + err.Error(),
		})
	}

	// create auth token
	_secret := os.Getenv("JWT_SECRET_KEY")
	month := (time.Hour * 24) * 30
	claims := jtoken.MapClaims{
		"Id":    admin.Id,
		"email": admin.Email,
		"role":  "admin",
		"exp":   time.Now().Add(month * 6).Unix(),
	}
	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)
	_token, err := token.SignedString([]byte(_secret))
	if err != nil {
		return c.Status(400).JSON(providerAuth.LoginProviderResDto{
			Status:  false,
			Message: "Token is not valid" + err.Error(),
		})
	}

	responseData := adminAuth.LoginVerifyOtpResDto{
		Status:  true,
		Message: "OTP verified successfully and successfully logged in",
		Data: adminAuth.GetAdminRes{
			Id:        admin.Id,
			FirstName: admin.FirstName,
			LastName:  admin.LastName,
			Email:     admin.Email,
			Image:     admin.Image,
			CreatedAt: admin.CreatedAt,
			UpdatedAt: admin.UpdatedAt,
		},
		Token: _token,
	}

	return c.Status(200).JSON(responseData)
}
