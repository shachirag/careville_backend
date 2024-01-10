package providerAuthenticate

import (
	"careville_backend/database"
	providerAuth "careville_backend/dto/provider/providerAuth"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Fetch provider By ID
// @Description Fetch provider By ID
// @Tags provider authorization
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "provider ID"
// @Produce json
// @Success 200 {object} providerAuth.GetProviderResDto
// @Router /provider/get-provider-info/{id} [get]
func FetchProviderById(c *fiber.Ctx) error {

	var provider entity.ServiceEntity

	// Get the user ID from the URL parameter
	userId := c.Params("id")
	objId, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return c.Status(400).JSON(providerAuth.GetProviderResDto{
			Status:  false,
			Message: "invalid objectId " + err.Error(),
		})

	}

	serviceColl := database.GetCollection("service")

	err = serviceColl.FindOne(ctx, bson.M{"_id": objId}).Decode(&provider)
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

	if provider.Role == "healthFacility" && provider.FacilityOrProfession == "hospClinic" {
		image = provider.HospClinic.Information.Image
		additionalDetails = provider.HospClinic.Information.AdditionalText
		isEmergencyAvailable = provider.HospClinic.Information.IsEmergencyAvailable
		address = providerAuth.Address(provider.HospClinic.Information.Address)
	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "laboratory" {
		image = provider.Laboratory.Information.Image
		additionalDetails = provider.Laboratory.Information.AdditionalText
		isEmergencyAvailable = provider.Laboratory.Information.IsEmergencyAvailable
		address = providerAuth.Address(provider.Laboratory.Information.Address)
	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "fitnessCenter" {
		image = provider.FitnessCenter.Information.Image
		additionalDetails = provider.FitnessCenter.Information.AdditionalText
		isEmergencyAvailable = provider.FitnessCenter.Information.IsEmergencyAvailable
		address = providerAuth.Address(provider.FitnessCenter.Information.Address)
	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "pharmacy" {
		image = provider.Pharmacy.Information.Image
		additionalDetails = provider.Pharmacy.Information.AdditionalText
		isEmergencyAvailable = provider.Pharmacy.Information.IsEmergencyAvailable
		address = providerAuth.Address(provider.Pharmacy.Information.Address)
	} else if provider.Role == "healthProfessional" && provider.FacilityOrProfession == "medicalLabScientist" {
		image = provider.MedicalLabScientist.Information.Image
		additionalDetails = provider.MedicalLabScientist.Information.AdditionalText
		isEmergencyAvailable = provider.MedicalLabScientist.Information.IsEmergencyAvailable
		address = providerAuth.Address(provider.MedicalLabScientist.Information.Address)
	} else if provider.Role == "healthProfessional" && provider.FacilityOrProfession == "nurse" {
		image = provider.Nurse.Information.Image
		additionalDetails = provider.Nurse.Information.AdditionalText
		isEmergencyAvailable = provider.Nurse.Information.IsEmergencyAvailable
		address = providerAuth.Address(provider.Nurse.Information.Address)
	} else if provider.Role == "healthProfessional" && provider.FacilityOrProfession == "doctor" {
		image = provider.Doctor.Information.Image
		additionalDetails = provider.Doctor.Information.AdditionalText
		isEmergencyAvailable = provider.Doctor.Information.IsEmergencyAvailable
		address = providerAuth.Address(provider.Doctor.Information.Address)
	} else if provider.Role == "healthProfessional" && provider.FacilityOrProfession == "physiotherapist" {
		image = provider.Physiotherapist.Information.Image
		additionalDetails = provider.Physiotherapist.Information.AdditionalText
		isEmergencyAvailable = provider.Physiotherapist.Information.IsEmergencyAvailable
		address = providerAuth.Address(provider.Physiotherapist.Information.Address)
	}

	providerRes := providerAuth.ProviderData{
		Id:                   provider.Id,
		Name:                 provider.Name,
		Email:                provider.Email,
		Image:                image,
		AdditionalDetails:    additionalDetails,
		Address:              address,
		IsEmergencyAvailable: isEmergencyAvailable,
		PhoneNumber:          providerAuth.PhoneNumber(provider.PhoneNumber),
		CreatedAt:            provider.CreatedAt,
		UpdatedAt:            provider.UpdatedAt,
		Notification: providerAuth.Notification{
			DeviceToken: provider.Notification.DeviceToken,
			DeviceType:  provider.Notification.DeviceType,
			IsEnabled:   provider.Notification.IsEnabled,
		},
	}

	return c.Status(fiber.StatusOK).JSON(providerAuth.GetProviderResDto{
		Status:   true,
		Message:  "provider data retrieved successfully",
		Provider: providerRes,
	})
}
