package doctorProfession

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	providerAuth "careville_backend/dto/provider/providerAuth"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Fetch personal information By ID
// @Description Fetch personal information By ID
// @Tags doctorProfession
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} providerAuth.GetProviderResDto
// @Router /provider/services/get-doctorProfession-personal-info [get]
func FetchDoctorProfessionPersonalDetailsById(c *fiber.Ctx) error {

	var provider entity.ServiceEntity

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	serviceColl := database.GetCollection("service")

	projection := bson.M{
		"role":                                    1,
		"facilityOrProfession":                    1,
		"serviceStatus":                           1,
		"doctor.information.image":                1,
		"doctor.information.name":                 1,
		"doctor.information.additionalText":       1,
		"doctor.information.isEmergencyAvailable": 1,
		"doctor.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
		"doctor.personalIdentificationDocs.nimc":    1,
		"doctor.personalIdentificationDocs.license": 1,
		"user.id":                       1,
		"user.firstName":                1,
		"user.lastName":                 1,
		"user.phoneNumber.dialCode":     1,
		"user.phoneNumber.countryCode":  1,
		"user.phoneNumber.number":       1,
		"user.notification.deviceType":  1,
		"user.notification.deviceToken": 1,
		"user.notification.isEnabled":   1,
		"user.email":                    1,
		"createdAt":                     1,
		"updatedAt":                     1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	err := serviceColl.FindOne(ctx, bson.M{"_id": providerData.ProviderId}, findOptions).Decode(&provider)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(providerAuth.GetProviderResDto{
				Status:  false,
				Message: "service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.GetProviderResDto{
			Status:  false,
			Message: "Failed to fetch service from MongoDB: " + err.Error(),
		})
	}

	var image string
	var additionalDetails string
	var isEmergencyAvailable bool
	var address providerAuth.Address
	var license string
	var nimc string
	var name string

	if provider.Doctor != nil {
		image = provider.Doctor.Information.Image
		additionalDetails = provider.Doctor.Information.AdditionalText
		isEmergencyAvailable = provider.Doctor.Information.IsEmergencyAvailable
		address = providerAuth.Address(provider.Doctor.Information.Address)
		license = provider.Doctor.PersonalIdentificationDocs.License
		nimc = provider.Doctor.PersonalIdentificationDocs.Nimc
		name = provider.Doctor.Information.Name
	}

	providerRes := providerAuth.ProviderResDto{
		User: providerAuth.UserData{
			Role: providerAuth.Role{
				Role:                 provider.Role,
				FacilityOrProfession: provider.FacilityOrProfession,
				ServiceStatus:        provider.ServiceStatus,
				Image:                image,
				Name:                 name,
				IsEmergencyAvailable: isEmergencyAvailable,
			},
			User: providerAuth.User{
				Id:          provider.Id,
				FirstName:   provider.User.FirstName,
				LastName:    provider.User.LastName,
				Email:       provider.User.Email,
				PhoneNumber: providerAuth.PhoneNumber(provider.User.PhoneNumber),
				Notification: providerAuth.Notification{
					DeviceToken: provider.User.Notification.DeviceToken,
					DeviceType:  provider.User.Notification.DeviceType,
					IsEnabled:   provider.User.Notification.IsEnabled,
				},
				CreatedAt: provider.CreatedAt,
				UpdatedAt: provider.UpdatedAt,
			},
		},
		AdditionalInformation: providerAuth.AdditionalInformation{
			AdditionalDetails: additionalDetails,
			Address:           address,
			Documents: providerAuth.Documents{
				Certificate: nimc,
				License:     license,
			},
		},
	}

	return c.Status(fiber.StatusOK).JSON(providerAuth.GetProviderResDto{
		Status:   true,
		Message:  "provider data retrieved successfully",
		Provider: providerRes,
	})
}
