package providerAuthenticate

import (
	"careville_backend/database"
	providerAuth "careville_backend/dto/provider/providerAuth"
	"careville_backend/entity"
	"os"
	"strings"
	"time"

	jtoken "github.com/golang-jwt/jwt/v4"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Login provider
// @Description Login provider
// @Tags provider authorization
// @Accept application/json
// @Param provider body providerAuth.LoginProviderReqDto true "login provider"
// @Produce json
// @Success 200 {object} providerAuth.LoginProviderResDto
// @Router /provider/login [post]
func LoginProvider(c *fiber.Ctx) error {

	var (
		providerColl = database.GetCollection("provider")
		data         providerAuth.LoginProviderReqDto
		provider     entity.ProviderEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(providerAuth.LoginProviderResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Find the user with email address from client
	err = providerColl.FindOne(ctx, bson.M{"email": data.Email}).Decode(&provider)
	if err != nil {
		// Check if there is no documents found error
		if err == mongo.ErrNoDocuments {
			return c.Status(400).JSON(providerAuth.LoginProviderResDto{
				Status:  false,
				Message: "Invalid credentials",
			})
		}

		return c.Status(500).JSON(providerAuth.LoginProviderResDto{
			Status:  false,
			Message: "Internal server error, while getting the provider: " + err.Error(),
		})
	}

	// Checking if passwords match
	err = bcrypt.CompareHashAndPassword([]byte(provider.Password), []byte(strings.TrimSpace(data.Password)))
	if err != nil {
		return c.Status(400).JSON(providerAuth.LoginProviderResDto{
			Status:  false,
			Message: "Invalid credentials",
		})
	}

	// create auth token
	_secret := os.Getenv("JWT_SECRET_KEY")
	// _token_exp := os.Getenv("JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT")
	// t, err := utils.CreateToken(user, _secret)
	month := (time.Hour * 24) * 30
	// Create the JWT claims, which includes the user ID and expiry time
	claims := jtoken.MapClaims{
		"Id":    provider.Id,
		"email": provider.Email,
		"role":  "provider",
		"exp":   time.Now().Add(month * 6).Unix(),
	}
	// Create token
	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	_token, err := token.SignedString([]byte(_secret))
	if err != nil {
		return c.Status(400).JSON(providerAuth.LoginProviderResDto{
			Status:  false,
			Message: "Token is not valid" + err.Error(),
		})
	}

	return c.Status(200).JSON(providerAuth.LoginProviderResDto{
		Status:  true,
		Message: "Successfully logged in.",
		Provider: providerAuth.ProviderResDto{
			Id:          provider.Id,
			Name:        provider.Name,
			Email:       provider.Email,
			Image:       provider.Image,
			PhoneNumber: providerAuth.PhoneNumber(provider.PhoneNumber),
			CreatedAt:   provider.CreatedAt,
			UpdatedAt:   provider.UpdatedAt,
		},
		Token: _token,
	})
}
