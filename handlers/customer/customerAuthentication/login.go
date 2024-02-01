package customerAuth

import (
	"careville_backend/database"
	customerAuth "careville_backend/dto/customer/customerAuth"
	"careville_backend/entity"

	"context"
	"os"
	"strings"
	"time"

	jtoken "github.com/golang-jwt/jwt/v4"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var ctx = context.Background()

// @Summary Login customer
// @Description Login customer
// @Tags customer authorization
// @Accept application/json
// @Param customer body customerAuth.LoginCustomerReqDto true "login customer"
// @Produce json
// @Success 200 {object} customerAuth.LoginCustomerResDto
// @Router /customer/login [post]
func LoginCustomer(c *fiber.Ctx) error {
	var (
		customerColl = database.GetCollection("customer")
		data         customerAuth.LoginCustomerReqDto
		customer     entity.CustomerEntity
	)

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(customerAuth.LoginCustomerResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	projection := bson.M{
		"_id":                        1,
		"firstName":                  1,
		"lastName":                   1,
		"image":                      1,
		"age":                        1,
		"sex":                        1,
		"familyMembers.id":           1,
		"familyMembers.name":         1,
		"familyMembers.age":          1,
		"familyMembers.sex":          1,
		"familyMembers.relationShip": 1,
		"address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
		"password":                 1,
		"phoneNumber.dialCode":     1,
		"phoneNumber.countryCode":  1,
		"phoneNumber.number":       1,
		"notification.deviceType":  1,
		"notification.deviceToken": 1,
		"notification.isEnabled":   1,
		"email":                    1,
		"createdAt":                1,
		"updatedAt":                1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	err = customerColl.FindOne(ctx, bson.M{"email": data.Email}, findOptions).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(customerAuth.LoginCustomerResDto{
				Status:  false,
				Message: "Invalid credentials",
			})
		}

		return c.Status(500).JSON(customerAuth.LoginCustomerResDto{
			Status:  false,
			Message: "Internal server error, while getting the user: " + err.Error(),
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(strings.TrimSpace(data.Password)))
	if err != nil {
		return c.Status(400).JSON(customerAuth.LoginCustomerResDto{
			Status:  false,
			Message: "Invalid credentials",
		})
	}

	_secret := os.Getenv("JWT_SECRET_KEY")
	month := (time.Hour * 24) * 30
	claims := jtoken.MapClaims{
		"Id":    customer.Id,
		"email": customer.Email,
		"role":  "customer",
		"type":  customer.Type,
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
		Message: "Successfully logged in.",
		Data: customerAuth.CustomerResDto{
			Id:        customer.Id,
			FirstName: customer.FirstName,
			LastName:  customer.LastName,
			Email:     customer.Email,
			Image:     customer.Image,
			PhoneNumber: customerAuth.PhoneNumber{
				DialCode:    customer.PhoneNumber.DialCode,
				CountryCode: customer.PhoneNumber.CountryCode,
				Number:      customer.PhoneNumber.Number,
			},
			Notification: customerAuth.Notification{
				DeviceToken: customer.Notification.DeviceToken,
				DeviceType:  customer.Notification.DeviceType,
				IsEnabled:   customer.Notification.IsEnabled,
			},
			Address: customerAuth.Address{
				Coordinates: customer.Address.Coordinates,
				Type:        customer.Address.Type,
				Add:         customer.Address.Add,
			},
			FamilyMembers: familyData,
			Age:           customer.Age,
			Sex:           customer.Sex,
			CreatedAt:     customer.CreatedAt,
			UpdatedAt:     customer.UpdatedAt,
		},
		Token: _token,
	})
}
