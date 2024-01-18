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
// @Param id path string true "provider ID"
// @Param provider formData providerAuth.UpdateProviderReqDto true "Update data of provider"
// @Param newProviderImage formData file false "provider profile image"
// @Produce json
// @Success 200 {object} providerAuth.UpdateProviderResDto
// @Router /provider/profile/update-provider-data [put]
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

	if provider.Role == "healthFacility" && provider.FacilityOrProfession == "hospClinic" {
		update = bson.M{"$set": bson.M{
			"firstName": data.FirstName,
			"lastName":  data.LastName,
			"phoneNumber": bson.M{
				"dialCode":    data.DialCode,
				"number":      data.PhoneNumber,
				"countryCode": data.CountryCode,
			},
			"updatedAt":                             time.Now().UTC(),
			"hospClinic.information.firstName":      data.FirstName,
			"hospClinic.information.lastName":       data.LastName,
			"hospClinic.information.additionalText": data.AdditionalText,
			"hospClinic.information.address": providerAuth.Address{
				Coordinates: []float64{longitude, latitude},
				Type:        "Point",
				Add:         data.Address,
			},
		},
		}
	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "laboratory" {
		update = bson.M{"$set": bson.M{
			"firstName": data.FirstName,
			"lastName":  data.LastName,
			"phoneNumber": bson.M{
				"dialCode":    data.DialCode,
				"number":      data.PhoneNumber,
				"countryCode": data.CountryCode,
			},
			"updatedAt":                             time.Now().UTC(),
			"laboratory.information.firstName":      data.FirstName,
			"laboratory.information.lastName":       data.LastName,
			"laboratory.information.additionalText": data.AdditionalText,
			"laboratory.information.address": providerAuth.Address{
				Coordinates: []float64{longitude, latitude},
				Type:        "Point",
				Add:         data.Address,
			},
		},
		}
	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "fitnessCenter" {
		update = bson.M{"$set": bson.M{
			"firstName": data.FirstName,
			"lastName":  data.LastName,
			"phoneNumber": bson.M{
				"dialCode":    data.DialCode,
				"number":      data.PhoneNumber,
				"countryCode": data.CountryCode,
			},
			"updatedAt":                                time.Now().UTC(),
			"fitnessCenter.information.firstName":      data.FirstName,
			"fitnessCenter.information.lastName":       data.LastName,
			"fitnessCenter.information.additionalText": data.AdditionalText,
			"fitnessCenter.information.address": providerAuth.Address{
				Coordinates: []float64{longitude, latitude},
				Type:        "Point",
				Add:         data.Address,
			},
		},
		}
	} else if provider.Role == "healthFacility" && provider.FacilityOrProfession == "pharmacy" {
		update = bson.M{"$set": bson.M{
			"firstName": data.FirstName,
			"lastName":  data.LastName,
			"phoneNumber": bson.M{
				"dialCode":    data.DialCode,
				"number":      data.PhoneNumber,
				"countryCode": data.CountryCode,
			},
			"updatedAt":                           time.Now().UTC(),
			"pharmacy.information.firstName":      data.FirstName,
			"pharmacy.information.lastName":       data.LastName,
			"pharmacy.information.additionalText": data.AdditionalText,
			"pharmacy.information.address": providerAuth.Address{
				Coordinates: []float64{longitude, latitude},
				Type:        "Point",
				Add:         data.Address,
			},
		},
		}
	} else if provider.Role == "healthProfessional" && provider.FacilityOrProfession == "medicalLabScientist" {
		update = bson.M{"$set": bson.M{
			"firstName": data.FirstName,
			"lastName":  data.LastName,
			"phoneNumber": bson.M{
				"dialCode":    data.DialCode,
				"number":      data.PhoneNumber,
				"countryCode": data.CountryCode,
			},
			"updatedAt": time.Now().UTC(),
			"medicalLabScientist.information.firstName":      data.FirstName,
			"medicalLabScientist.information.lastName":       data.LastName,
			"medicalLabScientist.information.additionalText": data.AdditionalText,
			"medicalLabScientist.information.address": providerAuth.Address{
				Coordinates: []float64{longitude, latitude},
				Type:        "Point",
				Add:         data.Address,
			},
		},
		}
	} else if provider.Role == "healthProfessional" && provider.FacilityOrProfession == "nurse" {
		update = bson.M{"$set": bson.M{
			"firstName": data.FirstName,
			"lastName":  data.LastName,
			"phoneNumber": bson.M{
				"dialCode":    data.DialCode,
				"number":      data.PhoneNumber,
				"countryCode": data.CountryCode,
			},
			"updatedAt":                        time.Now().UTC(),
			"nurse.information.firstName":      data.FirstName,
			"nurse.information.lastName":       data.LastName,
			"nurse.information.additionalText": data.AdditionalText,
			"nurse.information.address": providerAuth.Address{
				Coordinates: []float64{longitude, latitude},
				Type:        "Point",
				Add:         data.Address,
			},
		},
		}
	} else if provider.Role == "healthProfessional" && provider.FacilityOrProfession == "doctor" {
		update = bson.M{"$set": bson.M{
			"firstName": data.FirstName,
			"lastName":  data.LastName,
			"phoneNumber": bson.M{
				"dialCode":    data.DialCode,
				"number":      data.PhoneNumber,
				"countryCode": data.CountryCode,
			},
			"updatedAt":                         time.Now().UTC(),
			"doctor.information.firstName":      data.FirstName,
			"doctor.information.lastName":       data.LastName,
			"doctor.information.additionalText": data.AdditionalText,
			"doctor.information.address": providerAuth.Address{
				Coordinates: []float64{longitude, latitude},
				Type:        "Point",
				Add:         data.Address,
			},
		},
		}
	} else if provider.Role == "healthProfessional" && provider.FacilityOrProfession == "physiotherapist" {
		update = bson.M{"$set": bson.M{
			"firstName": data.FirstName,
			"lastName":  data.LastName,
			"phoneNumber": bson.M{
				"dialCode":    data.DialCode,
				"number":      data.PhoneNumber,
				"countryCode": data.CountryCode,
			},
			"updatedAt":                                  time.Now().UTC(),
			"physiotherapist.information.firstName":      data.FirstName,
			"physiotherapist.information.lastName":       data.LastName,
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
	})
}
