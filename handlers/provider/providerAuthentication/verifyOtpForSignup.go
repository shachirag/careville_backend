package providerAuthenticate

import (
	"careville_backend/database"
	providerAuth "careville_backend/dto/provider/providerAuth"
	"careville_backend/entity"
	"os"
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
// @Tags provider authorization
// @Accept application/json
// @Param provider body providerAuth.ProviderSignupVerifyOtpReqDto true "Verify 6 digit OTP and insert data into database"
// @Produce json
// @Success 200 {object} providerAuth.LoginProviderResDto
// @Router /provider/verify-otp-for-signup [post]
func VerifyOtpForSignup(c *fiber.Ctx) error {
	var (
		otpColl     = database.GetCollection("otp")
		serviceColl = database.GetCollection("service")
		data        providerAuth.ProviderSignupVerifyOtpReqDto
		otpData     entity.OtpEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(providerAuth.LoginProviderResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Error handling
	if data.EnteredOTP == "" {
		return c.Status(400).JSON(providerAuth.LoginProviderResDto{
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
			return c.Status(400).JSON(providerAuth.LoginProviderResDto{
				Status:  false,
				Message: "Invalid OTP",
			})
		}

		return c.Status(500).JSON(providerAuth.LoginProviderResDto{
			Status:  false,
			Message: "Internal server error, while getting the provider: " + err.Error(),
		})
	}

	// Compare the entered OTP with the OTP from the database
	if data.EnteredOTP != otpData.Otp {

		return c.Status(400).JSON(providerAuth.LoginProviderResDto{
			Status:  false,
			Message: "Invalid OTP",
		})
	}

	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.LoginProviderResDto{
			Status:  false,
			Message: "Failed to hash the password: " + err.Error(),
		})
	}

	filter := bson.M{
		"user.email": strings.ToLower(data.Email),
	}

	exists, err := serviceColl.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(500).JSON(providerAuth.LoginProviderResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	if exists > 0 {
		return c.Status(400).JSON(providerAuth.LoginProviderResDto{
			Status:  false,
			Message: "Email is already in use",
		})
	}

	// Create the customer document to be inserted into MongoDB
	id := primitive.NewObjectID()

	// Now that OTP is verified, proceed to insert the data into the database
	provider := entity.ServiceEntity{
		Id: id,
		User: entity.ProviderUser{
			FirstName: data.FirstName,
			LastName:  data.LastName,
			Email:     smallEmail,
			PhoneNumber: entity.PhoneNumber{
				DialCode:    data.DialCode,
				Number:      data.PhoneNumber,
				CountryCode: data.CountryCode,
			},
			Notification: entity.Notification{
				DeviceToken: data.DeviceToken,
				DeviceType:  data.DeviceType,
				IsEnabled:   true,
			},
			Password: string(hashedPassword),
		},
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	_, err = serviceColl.InsertOne(ctx, provider)
	if err != nil {
		return c.Status(500).JSON(providerAuth.LoginProviderResDto{
			Status:  false,
			Message: "Failed to insert provider data: " + err.Error(),
		})
	}

	// create auth token
	_secret := os.Getenv("JWT_SECRET_KEY")
	// _token_exp := os.Getenv("JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT")
	// t, err := utils.CreateToken(user, _secret)
	month := (time.Hour * 24) * 30
	// Create the JWT claims, which includes the user ID and expiry time
	claims := jtoken.MapClaims{
		"Id":                   provider.Id,
		"email":                provider.User.Email,
		"role":                 "provider",
		"serviceRole":          provider.Role,
		"facilityOrProfession": provider.FacilityOrProfession,
		"exp":                  time.Now().Add(month * 6).Unix(),
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

	var image string
	var name string
	var isEmergencyAvailable bool

	if provider.Role == "healthFacility" && provider.FacilityOrProfession == "hospClinic" {
		image = provider.HospClinic.Information.Image
		name = provider.HospClinic.Information.Name
		isEmergencyAvailable = provider.HospClinic.Information.IsEmergencyAvailable
	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "laboratory" {
		image = provider.Laboratory.Information.Image
		name = provider.Laboratory.Information.Name
		isEmergencyAvailable = provider.Laboratory.Information.IsEmergencyAvailable
	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "fitnessCenter" {
		image = provider.FitnessCenter.Information.Image
		name = provider.FitnessCenter.Information.Name
		isEmergencyAvailable = provider.FitnessCenter.Information.IsEmergencyAvailable
	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "pharmacy" {
		image = provider.Pharmacy.Information.Image
		name = provider.Pharmacy.Information.Name
		isEmergencyAvailable = provider.Pharmacy.Information.IsEmergencyAvailable
	} else if provider.Role == "healthProfessional" && provider.FacilityOrProfession == "medicalLabScientist" {
		image = provider.MedicalLabScientist.Information.Image
		name = provider.MedicalLabScientist.Information.Name
		isEmergencyAvailable = provider.MedicalLabScientist.Information.IsEmergencyAvailable
	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "nurse" {
		image = provider.Nurse.Information.Image
		name = provider.Nurse.Information.Name
		isEmergencyAvailable = provider.Nurse.Information.IsEmergencyAvailable
	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "doctor" {
		image = provider.Doctor.Information.Image
		name = provider.Doctor.Information.Name
		isEmergencyAvailable = provider.Doctor.Information.IsEmergencyAvailable
	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "physiotherapist" {
		image = provider.Physiotherapist.Information.Image
		name = provider.Physiotherapist.Information.Name
		isEmergencyAvailable = provider.Physiotherapist.Information.IsEmergencyAvailable
	}

	role := providerAuth.Role{}
	if err == nil {
		role = providerAuth.Role{
			Role:                 provider.Role,
			FacilityOrProfession: provider.FacilityOrProfession,
			ServiceStatus:        provider.ServiceStatus,
			Image:                image,
			Name:                 name,
			IsEmergencyAvailable: isEmergencyAvailable,
		}
	}

	return c.Status(200).JSON(providerAuth.LoginProviderResDto{
		Status:  true,
		Message: "OTP verified successfully",
		Provider: providerAuth.ProviderRespDto{
			Role: role,
			User: providerAuth.User{
				Id:        provider.Id,
				FirstName: provider.User.FirstName,
				LastName:  provider.User.LastName,
				Email:     provider.User.Email,
				PhoneNumber: providerAuth.PhoneNumber{
					DialCode:    provider.User.PhoneNumber.DialCode,
					Number:      provider.User.PhoneNumber.Number,
					CountryCode: provider.User.PhoneNumber.CountryCode,
				},
				Notification: providerAuth.Notification{
					DeviceToken: provider.User.Notification.DeviceToken,
					DeviceType:  provider.User.Notification.DeviceType,
					IsEnabled:   provider.User.Notification.IsEnabled,
				},
				CreatedAt: provider.CreatedAt,
				UpdatedAt: provider.UpdatedAt,
			},
		},
		Token: _token,
	})
}
