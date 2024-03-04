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
	"go.mongodb.org/mongo-driver/mongo/options"
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
				Message: "Provider not found",
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

		if provider.HospClinic != nil {
			image = provider.HospClinic.Information.Image
			name = provider.HospClinic.Information.Name
			isEmergencyAvailable = provider.HospClinic.Information.IsEmergencyAvailable
		}

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

		if provider.Laboratory != nil {
			image = provider.Laboratory.Information.Image
			name = provider.Laboratory.Information.Name
			isEmergencyAvailable = provider.Laboratory.Information.IsEmergencyAvailable
		}

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

		if provider.FitnessCenter != nil {
			image = provider.FitnessCenter.Information.Image
			name = provider.FitnessCenter.Information.Name
			isEmergencyAvailable = provider.FitnessCenter.Information.IsEmergencyAvailable
		}

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

		if provider.Pharmacy != nil {
			image = provider.Pharmacy.Information.Image
			name = provider.Pharmacy.Information.Name
			isEmergencyAvailable = provider.Pharmacy.Information.IsEmergencyAvailable
		}

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

		if provider.MedicalLabScientist != nil {
			image = provider.MedicalLabScientist.Information.Image
			name = provider.MedicalLabScientist.Information.Name
			isEmergencyAvailable = provider.MedicalLabScientist.Information.IsEmergencyAvailable
		}

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

		if provider.Nurse != nil {
			image = provider.Nurse.Information.Image
			name = provider.Nurse.Information.Name
			isEmergencyAvailable = provider.Nurse.Information.IsEmergencyAvailable
		}

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

		if provider.Doctor != nil {
			image = provider.Doctor.Information.Image
			name = provider.Doctor.Information.Name
			isEmergencyAvailable = provider.Doctor.Information.IsEmergencyAvailable
		}

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

		if provider.Physiotherapist != nil {
			image = provider.Physiotherapist.Information.Image
			name = provider.Physiotherapist.Information.Name
			isEmergencyAvailable = provider.Physiotherapist.Information.IsEmergencyAvailable
		}

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
			"updatedAt": time.Now().UTC(),
		},
		}
	}

	opts := options.Update().SetUpsert(true)
	session, err := database.GetMongoClient().StartSession()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.UpdateProviderResDto{
			Status:  false,
			Message: "Failed to start session",
		})
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		updateRes, err := serviceColl.UpdateOne(sessCtx, filter, update, opts)
		if err != nil {
			return nil, err
		}

		if updateRes.MatchedCount == 0 {
			return nil, mongo.ErrNoDocuments
		}

		filter := bson.M{"serviceId": providerData.ProviderId}

		var appointment entity.AppointmentEntity
		err = database.GetCollection("appointment").FindOne(ctx, filter).Decode(&appointment)
		if err != nil {
			return nil, err
		}

		var appointmentUpdate bson.M
		switch appointment.Role {
		case "healthFacility":
			switch appointment.FacilityOrProfession {
			case "hospital":
				appointmentUpdate = bson.M{"$set": bson.M{
					"hospital.information.address": providerAuth.Address{
						Coordinates: []float64{longitude, latitude},
						Type:        "Point",
						Add:         data.Address,
					},
				}}
			case "laboratory":
				appointmentUpdate = bson.M{"$set": bson.M{
					"laboratory.information.address": providerAuth.Address{
						Coordinates: []float64{longitude, latitude},
						Type:        "Point",
						Add:         data.Address,
					},
				}}
			case "fitnessCenter":
				appointmentUpdate = bson.M{"$set": bson.M{
					"fitnessCenter.information.address": providerAuth.Address{
						Coordinates: []float64{longitude, latitude},
						Type:        "Point",
						Add:         data.Address,
					},
				}}
			case "pharmacy":
				appointmentUpdate = bson.M{"$set": bson.M{
					"pharmacy.information.address": providerAuth.Address{
						Coordinates: []float64{longitude, latitude},
						Type:        "Point",
						Add:         data.Address,
					},
				}}
			}
		case "healthProfessional":
			switch appointment.FacilityOrProfession {
			case "doctor":
				appointmentUpdate = bson.M{"$set": bson.M{
					"doctor.information.address": providerAuth.Address{
						Coordinates: []float64{longitude, latitude},
						Type:        "Point",
						Add:         data.Address,
					},
				}}
			case "physiotherapist":
				appointmentUpdate = bson.M{"$set": bson.M{
					"physiotherapist.information.address": providerAuth.Address{
						Coordinates: []float64{longitude, latitude},
						Type:        "Point",
						Add:         data.Address,
					},
				}}
			case "medicalLabScientist":
				appointmentUpdate = bson.M{"$set": bson.M{
					"medicalLabScientist.information.address": providerAuth.Address{
						Coordinates: []float64{longitude, latitude},
						Type:        "Point",
						Add:         data.Address,
					},
				}}
			case "nurse":
				appointmentUpdate = bson.M{"$set": bson.M{
					"nurse.information.address": providerAuth.Address{
						Coordinates: []float64{longitude, latitude},
						Type:        "Point",
						Add:         data.Address,
					},
				}}
			}
		}

		_, err = database.GetCollection("appointment").UpdateMany(sessCtx, filter, appointmentUpdate)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.UpdateProviderResDto{
			Status:  false,
			Message: "Failed to update appointment data: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(providerAuth.UpdateProviderResDto{
		Status:  true,
		Message: "Provider data updated successfully",
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
