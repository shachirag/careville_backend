package providerAuthenticate

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	providerAuth "careville_backend/dto/provider/providerAuth"
	"careville_backend/entity"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Update provider
// @Description Update provider
// @Tags provider authorization
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param provider formData providerAuth.UpdateProviderReqDto true "Update data of provider"
// @Param newProviderImage formData file false "provider profile image"
// @Produce json
// @Success 200 {object} providerAuth.UpdateProviderResDto
// @Router /provider/profile/update-provider-info [put]
func UpdateProvider(c *fiber.Ctx) error {

	var (
		serviceColl = database.GetCollection("service")
		data        providerAuth.UpdateProviderReqDto
		provider    entity.ServiceEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(providerAuth.UpdateProviderResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	filter := bson.M{"_id": providerData.ProviderId}
	err = serviceColl.FindOne(ctx, filter).Decode(&provider)
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

	longitude, err := strconv.ParseFloat(data.Longitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(providerAuth.UpdateProviderResDto{
			Status:  false,
			Message: "Invalid longitude format",
		})
	}

	latitude, err := strconv.ParseFloat(data.Latitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(providerAuth.UpdateProviderResDto{
			Status:  false,
			Message: "Invalid latitude format",
		})
	}

	update := bson.M{}
	var image string
	var name string
	var isEmergencyAvailable bool

	if provider.Role == "healthFacility" && provider.FacilityOrProfession == "hospClinic" {
		image = provider.HospClinic.Information.Image
		name = provider.HospClinic.Information.Name
		isEmergencyAvailable = provider.HospClinic.Information.IsEmergencyAvailable
		update = bson.M{"$set": bson.M{
			"user.firstName": data.FirstName,
			"user.lastName":  data.LastName,
			"user.phoneNumber": bson.M{
				"dialCode":    data.DialCode,
				"number":      data.PhoneNumber,
				"countryCode": data.CountryCode,
			},
			"hospClinic.information.additionalText": data.AdditionalText,
			"hospClinic.information.address": providerAuth.Address{
				Coordinates: []float64{longitude, latitude},
				Type:        "Point",
				Add:         data.Address,
			},
			"updatedAt": time.Now().UTC(),
		},
		}
	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "laboratory" {
		image = provider.Laboratory.Information.Image
		name = provider.Laboratory.Information.Name
		isEmergencyAvailable = provider.HospClinic.Information.IsEmergencyAvailable
		update = bson.M{"$set": bson.M{
			"user.firstName": data.FirstName,
			"user.lastName":  data.LastName,
			"user.phoneNumber": bson.M{
				"dialCode":    data.DialCode,
				"number":      data.PhoneNumber,
				"countryCode": data.CountryCode,
			},
			"laboratory.information.additionalText": data.AdditionalText,
			"laboratory.information.address": providerAuth.Address{
				Coordinates: []float64{longitude, latitude},
				Type:        "Point",
				Add:         data.Address,
			},
			"updatedAt": time.Now().UTC(),
		},
		}
	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "fitnessCenter" {
		image = provider.FitnessCenter.Information.Image
		name = provider.FitnessCenter.Information.Name
		isEmergencyAvailable = provider.FitnessCenter.Information.IsEmergencyAvailable
		update = bson.M{"$set": bson.M{
			"user.firstName": data.FirstName,
			"user.lastName":  data.LastName,
			"user.phoneNumber": bson.M{
				"dialCode":    data.DialCode,
				"number":      data.PhoneNumber,
				"countryCode": data.CountryCode,
			},
			"fitnessCenter.information.additionalText": data.AdditionalText,

			"fitnessCenter.information.address": providerAuth.Address{
				Coordinates: []float64{longitude, latitude},
				Type:        "Point",
				Add:         data.Address,
			},
			"updatedAt": time.Now().UTC(),
		},
		}
	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "pharmacy" {
		image = provider.Pharmacy.Information.Image
		name = provider.Pharmacy.Information.Name
		isEmergencyAvailable = provider.Pharmacy.Information.IsEmergencyAvailable
		update = bson.M{"$set": bson.M{
			"user.firstName": data.FirstName,
			"user.lastName":  data.LastName,
			"user.phoneNumber": bson.M{
				"dialCode":    data.DialCode,
				"number":      data.PhoneNumber,
				"countryCode": data.CountryCode,
			},
			"pharmacy.information.additionalText": data.AdditionalText,
			"pharmacy.information.address": providerAuth.Address{
				Coordinates: []float64{longitude, latitude},
				Type:        "Point",
				Add:         data.Address,
			},
			"updatedAt": time.Now().UTC(),
		},
		}
	} else if provider.Role == "healthProfessional" && provider.FacilityOrProfession == "medicalLabScientist" {
		image = provider.MedicalLabScientist.Information.Image
		name = provider.MedicalLabScientist.Information.Name
		isEmergencyAvailable = provider.MedicalLabScientist.Information.IsEmergencyAvailable
		update = bson.M{"$set": bson.M{
			"user.firstName": data.FirstName,
			"user.lastName":  data.LastName,
			"user.phoneNumber": bson.M{
				"dialCode":    data.DialCode,
				"number":      data.PhoneNumber,
				"countryCode": data.CountryCode,
			},
			"medicalLabScientist.information.additionalText": data.AdditionalText,

			"medicalLabScientist.information.address": providerAuth.Address{
				Coordinates: []float64{longitude, latitude},
				Type:        "Point",
				Add:         data.Address,
			},
			"updatedAt": time.Now().UTC(),
		},
		}
	} else if provider.Role == "healthProfessional" && provider.FacilityOrProfession == "nurse" {
		image = provider.Nurse.Information.Image
		name = provider.Nurse.Information.Name
		isEmergencyAvailable = provider.Nurse.Information.IsEmergencyAvailable
		update = bson.M{"$set": bson.M{
			"user.firstName": data.FirstName,
			"user.lastName":  data.LastName,
			"user.phoneNumber": bson.M{
				"dialCode":    data.DialCode,
				"number":      data.PhoneNumber,
				"countryCode": data.CountryCode,
			},
			"nurse.information.additionalText": data.AdditionalText,

			"nurse.information.address": providerAuth.Address{
				Coordinates: []float64{longitude, latitude},
				Type:        "Point",
				Add:         data.Address,
			},
			"updatedAt": time.Now().UTC(),
		},
		}
	} else if provider.Role == "healthProfessional" && provider.FacilityOrProfession == "doctor" {
		image = provider.Doctor.Information.Image
		name = provider.Doctor.Information.Name
		isEmergencyAvailable = provider.Doctor.Information.IsEmergencyAvailable
		update = bson.M{"$set": bson.M{
			"user.firstName": data.FirstName,
			"user.lastName":  data.LastName,
			"user.phoneNumber": bson.M{
				"dialCode":    data.DialCode,
				"number":      data.PhoneNumber,
				"countryCode": data.CountryCode,
			},
			"doctor.information.additionalText": data.AdditionalText,

			"doctor.information.address": providerAuth.Address{
				Coordinates: []float64{longitude, latitude},
				Type:        "Point",
				Add:         data.Address,
			},
			"updatedAt": time.Now().UTC(),
		},
		}
	} else if provider.Role == "healthProfessional" && provider.FacilityOrProfession == "physiotherapist" {
		image = provider.Physiotherapist.Information.Image
		name = provider.Physiotherapist.Information.Name
		isEmergencyAvailable = provider.Physiotherapist.Information.IsEmergencyAvailable
		update = bson.M{"$set": bson.M{
			"user.firstName": data.FirstName,
			"user.lastName":  data.LastName,
			"user.phoneNumber": bson.M{
				"dialCode":    data.DialCode,
				"number":      data.PhoneNumber,
				"countryCode": data.CountryCode,
			},
			"physiotherapist.information.additionalText": data.AdditionalText,

			"physiotherapist.information.address": providerAuth.Address{
				Coordinates: []float64{longitude, latitude},
				Type:        "Point",
				Add:         data.Address,
			},
		},
		}
	}

	// Execute the update operation
	updateRes, err := serviceColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.UpdateProviderResDto{
			Status:  false,
			Message: "Failed to update provider data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(providerAuth.UpdateProviderResDto{
			Status:  false,
			Message: "provider not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(providerAuth.UpdateProviderResDto{
		Status:  true,
		Message: "provider data updated successfully",
		Data: providerAuth.UpdateProfileData{
			Role: providerAuth.Role{
				Role:                 provider.Role,
				FacilityOrProfession: provider.FacilityOrProfession,
				ServiceStatus:        provider.ServiceStatus,
				Name:                 name,
				Image:                image,
				IsEmergencyAvailable: isEmergencyAvailable,
			},
			User: providerAuth.User{
				Id:        providerData.ProviderId,
				FirstName: data.FirstName,
				LastName:  data.LastName,
				Email:     provider.User.Email,
				Notification: providerAuth.Notification{
					DeviceToken: provider.User.Notification.DeviceToken,
					DeviceType:  provider.User.Notification.DeviceType,
					IsEnabled:   provider.User.Notification.IsEnabled,
				},
				PhoneNumber: providerAuth.PhoneNumber{
					DialCode:    provider.User.PhoneNumber.DialCode,
					Number:      provider.User.PhoneNumber.Number,
					CountryCode: provider.User.PhoneNumber.CountryCode,
				},
				UpdatedAt: time.Now().UTC(),
				CreatedAt: provider.CreatedAt,
			},
		},
	})
}
