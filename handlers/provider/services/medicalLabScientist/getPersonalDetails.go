package medicalLabScientist

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
// @Tags medicalLabScientist
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} providerAuth.GetProviderResDto
// @Router /provider/services/get-medicalLabScientist-personal-details [get]
func FetchMedicalLabScientistPersonalDetailsById(c *fiber.Ctx) error {

	var provider entity.ServiceEntity

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	serviceColl := database.GetCollection("service")

	projection := bson.M{
		"role":                                                 1,
		"facilityOrProfession":                                 1,
		"serviceStatus":                                        1,
		"medicalLabScientist.information.image":                1,
		"medicalLabScientist.information.name":                 1,
		"medicalLabScientist.information.additionalText":       1,
		"medicalLabScientist.information.isEmergencyAvailable": 1,
		"medicalLabScientist.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
		"medicalLabScientist.personalIdentificationDocs.nimc":    1,
		"medicalLabScientist.personalIdentificationDocs.license": 1,
		"user.id":                      1,
		"user.firstName":               1,
		"user.lastName":                1,
		"user.phoneNumber.dialCode":    1,
		"user.phoneNumber.countryCode": 1,
		"user.phoneNumber.number":      1,
		"user.notification.deviceType":  1,
		"user.notification.deviceToken": 1,
		"user.notification.isEnabled":   1,
		"user.email":                   1,
		"createdAt":                    1,
		"updatedAt":                    1,
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

	if provider.MedicalLabScientist != nil {
		image = provider.MedicalLabScientist.Information.Image
		additionalDetails = provider.MedicalLabScientist.Information.AdditionalText
		isEmergencyAvailable = provider.MedicalLabScientist.Information.IsEmergencyAvailable
		address = providerAuth.Address(provider.MedicalLabScientist.Information.Address)
		license = provider.MedicalLabScientist.PersonalIdentificationDocs.License
		nimc = provider.MedicalLabScientist.PersonalIdentificationDocs.Nimc
		name = provider.MedicalLabScientist.Information.Name
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
