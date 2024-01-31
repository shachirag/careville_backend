package nurse

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
// @Tags nurse
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} providerAuth.GetProviderResDto
// @Router /provider/services/get-nurse-personal-info [get]
func FetchNursePersonalDetailsById(c *fiber.Ctx) error {

	var provider entity.ServiceEntity

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	serviceColl := database.GetCollection("service")

	projection := bson.M{
		"role":                                   1,
		"facilityOrProfession":                   1,
		"serviceStatus":                          1,
		"nurse.information.image":                1,
		"nurse.information.name":                 1,
		"nurse.information.additionalText":       1,
		"nurse.information.isEmergencyAvailable": 1,
		"nurse.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
		"nurse.personalIdentificationDocs.nimc":    1,
		"nurse.personalIdentificationDocs.license": 1,
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
				Message: "provider not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.GetProviderResDto{
			Status:  false,
			Message: "Failed to fetch provider from MongoDB: " + err.Error(),
		})
	}

	var image string
	var additionalDetails string
	var isEmergencyAvailable bool
	var address providerAuth.Address
	var license string
	var nimc string
	var name string

	if provider.Nurse != nil {
		image = provider.Nurse.Information.Image
		additionalDetails = provider.Nurse.Information.AdditionalText
		isEmergencyAvailable = provider.Nurse.Information.IsEmergencyAvailable
		address = providerAuth.Address(provider.Nurse.Information.Address)
		license = provider.Nurse.PersonalIdentificationDocs.License
		nimc = provider.Nurse.PersonalIdentificationDocs.Nimc
		name = provider.Nurse.Information.Name
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
