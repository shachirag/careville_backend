package providerAuthenticate

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	providerAuth "careville_backend/dto/provider/providerAuth"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Fetch provider By ID
// @Description Fetch provider By ID
// @Tags provider authorization
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} providerAuth.GetProviderResDto
// @Router /provider/profile/get-provider-info [get]
func FetchProviderById(c *fiber.Ctx) error {

	var provider entity.ServiceEntity

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	serviceColl := database.GetCollection("service")

	err := serviceColl.FindOne(ctx, bson.M{"_id": providerData.ProviderId}).Decode(&provider)
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
	var certificate string
	var name string

	if provider.Role == "healthFacility" && provider.FacilityOrProfession == "hospClinic" {

		if provider.HospClinic != nil {
			image = provider.HospClinic.Information.Image
			additionalDetails = provider.HospClinic.Information.AdditionalText
			isEmergencyAvailable = provider.HospClinic.Information.IsEmergencyAvailable
			address = providerAuth.Address(provider.HospClinic.Information.Address)
			license = provider.HospClinic.Documents.License
			certificate = provider.HospClinic.Documents.Certificate
			name = provider.HospClinic.Information.Name
		}

	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "laboratory" {

		if provider.Laboratory != nil {
			image = provider.Laboratory.Information.Image
			additionalDetails = provider.Laboratory.Information.AdditionalText
			isEmergencyAvailable = provider.Laboratory.Information.IsEmergencyAvailable
			address = providerAuth.Address(provider.Laboratory.Information.Address)
			license = provider.Laboratory.Documents.License
			certificate = provider.Laboratory.Documents.Certificate
			name = provider.Laboratory.Information.Name
		}

	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "fitnessCenter" {

		if provider.FitnessCenter != nil {
			image = provider.FitnessCenter.Information.Image
			additionalDetails = provider.FitnessCenter.Information.AdditionalText
			isEmergencyAvailable = provider.FitnessCenter.Information.IsEmergencyAvailable
			address = providerAuth.Address(provider.FitnessCenter.Information.Address)
			license = provider.FitnessCenter.Documents.License
			certificate = provider.FitnessCenter.Documents.Certificate
			name = provider.FitnessCenter.Information.Name
		}

	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "pharmacy" {

		if provider.Pharmacy != nil {
			image = provider.Pharmacy.Information.Image
			additionalDetails = provider.Pharmacy.Information.AdditionalText
			isEmergencyAvailable = provider.Pharmacy.Information.IsEmergencyAvailable
			address = providerAuth.Address(provider.Pharmacy.Information.Address)
			license = provider.Pharmacy.Documents.License
			certificate = provider.Pharmacy.Documents.Certificate
			name = provider.Pharmacy.Information.Name
		}

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
				Certificate: certificate,
				License:     license,
			},
		},
	}

	return c.Status(fiber.StatusOK).JSON(providerAuth.GetProviderResDto{
		Status:   true,
		Message:  "Provider data retrieved successfully",
		Provider: providerRes,
	})
}
