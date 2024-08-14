package customerAuth

import (
	"careville_backend/database"
	customerAuth "careville_backend/dto/customer/customerAuth"
	"careville_backend/entity"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	jtoken "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Verify OTP for signup
// @Description Verify the entered 6 digit OTP
// @Tags customer authorization
// @Accept application/json
// @Param customer body customerAuth.CustomerSignupVerifyOtpReqDto true "Verify 6 digit OTP and insert data into database"
// @Produce json
// @Success 200 {object} customerAuth.LoginCustomerResDto
// @Router /customer/verify-otp-for-signup [post]
func VerifyOtpForSignup(c *fiber.Ctx) error {
	var (
		otpColl      = database.GetCollection("otp")
		customerColl = database.GetCollection("customer")
		data         customerAuth.CustomerSignupVerifyOtpReqDto
		customer     entity.CustomerEntity
		otpData      entity.OtpEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(customerAuth.LoginCustomerResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Error handling
	if data.EnteredOTP == "" {
		return c.Status(400).JSON(customerAuth.LoginCustomerResDto{
			Status:  false,
			Message: "Entered OTP is required",
		})
	}

	smallEmail := strings.ToLower(data.Email)

	// Find the user with email address from client
	err = otpColl.FindOne(ctx, bson.M{"email": smallEmail}, options.FindOne().SetSort(bson.M{"createdAt": -1})).Decode(&otpData)
	if err != nil {
		// Check if there is no documents found error
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(customerAuth.LoginCustomerResDto{
				Status:  false,
				Message: "Invalid OTP",
			})
		}

		return c.Status(500).JSON(customerAuth.LoginCustomerResDto{
			Status:  false,
			Message: "Internal server error, while getting the provider: " + err.Error(),
		})
	}

	if data.EnteredOTP != otpData.Otp {

		return c.Status(400).JSON(customerAuth.LoginCustomerResDto{
			Status:  false,
			Message: "Invalid OTP",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(customerAuth.LoginCustomerResDto{
			Status:  false,
			Message: "Failed to hash the password: " + err.Error(),
		})
	}

	filter := bson.M{
		"email": strings.ToLower(data.Email),
	}

	exists, err := customerColl.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(500).JSON(customerAuth.LoginCustomerResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	if exists > 0 {
		return c.Status(400).JSON(customerAuth.LoginCustomerResDto{
			Status:  false,
			Message: "Email is already in use",
		})
	}

	longitude, err := strconv.ParseFloat(data.Longitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(customerAuth.LoginCustomerResDto{
			Status:  false,
			Message: "Invalid longitude format",
		})
	}

	latitude, err := strconv.ParseFloat(data.Latitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(customerAuth.LoginCustomerResDto{
			Status:  false,
			Message: "Invalid latitude format",
		})
	}

	id := primitive.NewObjectID()

	customer = entity.CustomerEntity{
		Id:        id,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		Email:     smallEmail,
		PhoneNumber: entity.PhoneNumber{
			DialCode:    data.PhoneNumber.DialCode,
			Number:      data.PhoneNumber.Number,
			CountryCode: data.PhoneNumber.CountryCode,
		},
		Notification: entity.Notification{
			DeviceToken: data.DeviceToken,
			DeviceType:  data.DeviceType,
			IsEnabled:   true,
		},
		IsDeleted: false,
		Address: entity.Address{
			Coordinates: []float64{longitude, latitude},
			Add:         data.Address,
			Type:        "Point",
		},
		Sex:       data.Sex,
		Age:       data.Age,
		Password:  string(hashedPassword),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	_, err = customerColl.InsertOne(ctx, customer)
	if err != nil {
		return c.Status(500).JSON(customerAuth.LoginCustomerResDto{
			Status:  false,
			Message: "Failed to insert provider data: " + err.Error(),
		})
	}

	_secret := os.Getenv("JWT_SECRET_KEY")
	month := (time.Hour * 24) * 30
	claims := jtoken.MapClaims{
		"Id":    customer.Id,
		"email": customer.Email,
		"role":  "customer",
		"exp":   time.Now().Add(month * 6).Unix(),
	}
	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)
	_token, err := token.SignedString([]byte(_secret))
	if err != nil {
		return c.Status(400).JSON(customerAuth.LoginCustomerResDto{
			Status:  false,
			Message: "Token is not valid" + err.Error(),
		})
	}

	familyData := make([]customerAuth.FamilyMembers, 0)
	if customer.FamilyMembers != nil {
		for _, family := range customer.FamilyMembers {
			familyData = append(familyData, customerAuth.FamilyMembers{
				Id:           family.Id,
				Name:         family.Name,
				RelationShip: family.RelationShip,
				Age:          family.Age,
				Sex:          family.Sex,
			})
		}
	}

	return c.Status(200).JSON(customerAuth.LoginCustomerResDto{
		Status:  true,
		Message: "OTP verified successfully",
		Data: customerAuth.CustomerResDto{
			Id:        customer.Id,
			FirstName: customer.FirstName,
			LastName:  customer.LastName,
			Email:     customer.Email,
			PhoneNumber: customerAuth.PhoneNumber{
				DialCode:    customer.PhoneNumber.DialCode,
				Number:      customer.PhoneNumber.Number,
				CountryCode: customer.PhoneNumber.CountryCode,
			},
			Notification: customerAuth.Notification{
				DeviceToken: customer.Notification.DeviceToken,
				DeviceType:  customer.Notification.DeviceType,
				IsEnabled:   customer.Notification.IsEnabled,
			},
			FamilyMembers: familyData,
			Address:       customerAuth.Address(customer.Address),
			Sex:           customer.Sex,
			Age:           customer.Age,
			CreatedAt:     customer.CreatedAt,
			UpdatedAt:     customer.UpdatedAt,
		},
		Token: _token,
	})
}
