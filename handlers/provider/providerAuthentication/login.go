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
		serviceColl = database.GetCollection("service")
		data        providerAuth.LoginProviderReqDto
		provider    entity.ServiceEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(providerAuth.LoginProviderResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	filter := bson.M{
		"email": strings.ToLower(data.Email),
	}

	// Find the user with email address from client
	err = serviceColl.FindOne(ctx, filter).Decode(&provider)
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
	month := (time.Hour * 24) * 30
	claims := jtoken.MapClaims{
		"Id":                   provider.Id,
		"email":                provider.Email,
		"role":                 "provider",
		"serviceRole":          provider.Role,
		"facilityOrProfession": provider.FacilityOrProfession,
		"exp":                  time.Now().Add(month * 6).Unix(),
	}
	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)
	_token, err := token.SignedString([]byte(_secret))
	if err != nil {
		return c.Status(400).JSON(providerAuth.LoginProviderResDto{
			Status:  false,
			Message: "Token is not valid" + err.Error(),
		})
	}

	role := providerAuth.Role{}
	if err == nil {
		role = providerAuth.Role{
			Role:                 provider.Role,
			FacilityOrProfession: provider.FacilityOrProfession,
			Status:               provider.Status,
		}
	}

	var image string

	if provider.Role == "healthFacility" && provider.FacilityOrProfession == "hospClinic" {
		image = provider.HospClinic.Information.Image
	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "laboratory" {
		image = provider.Laboratory.Information.Image
	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "fitnessCenter" {
		image = provider.FitnessCenter.Information.Image
	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "pharmacy" {
		image = provider.Pharmacy.Information.Image
	} else if provider.Role == "healthProfessional" && provider.FacilityOrProfession == "medicalLabScientist" {
		image = provider.MedicalLabScientist.Information.Image
	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "nurse" {
		image = provider.Nurse.Information.Image
	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "doctor" {
		image = provider.Doctor.Information.Image
	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "physiotherapist" {
		image = provider.Physiotherapist.Information.Image
	}

	return c.Status(200).JSON(providerAuth.LoginProviderResDto{
		Status:  true,
		Message: "Successfully logged in.",
		Provider: providerAuth.ProviderRespDto{
			Role: role,
			User: providerAuth.User{
				Id:        provider.Id,
				FirstName: provider.FirstName,
				LastName:  provider.LastName,
				Email:     provider.Email,
				Image:     image,
				PhoneNumber: providerAuth.PhoneNumber{
					DialCode:    provider.PhoneNumber.DialCode,
					Number:      provider.PhoneNumber.Number,
					CountryCode: provider.PhoneNumber.CountryCode,
				},
				Notification: providerAuth.Notification{
					DeviceToken: provider.Notification.DeviceToken,
					DeviceType:  provider.Notification.DeviceType,
					IsEnabled:   provider.Notification.IsEnabled,
				},
				CreatedAt: provider.CreatedAt,
				UpdatedAt: provider.UpdatedAt,
			},
		},
		Token: _token,
	})
}
