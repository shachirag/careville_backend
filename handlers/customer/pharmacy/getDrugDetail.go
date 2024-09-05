package pharmacy

import (
	"careville_backend/database"
	pharmacy "careville_backend/dto/customer/pharmacy"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get appointment by ID
// @Tags customer appointments
// @Description Get appointment by ID
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "appointment ID"
// @Produce json
// @Success 200 {object} pharmacy.GetPharmacyDrugsDetailResDto
// @Router /customer/healthFacility/appointment/pharmacy-drug/{id} [get]
func GetPharmacyAppointmentByID(c *fiber.Ctx) error {
	var (
		appointmentColl = database.GetCollection("appointment")
		serviceColl     = database.GetCollection("service")
	)

	idParam := c.Params("id")
	appointmentID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pharmacy.GetPharmacyDrugsDetailResDto{
			Status:  false,
			Message: "Invalid appointment ID",
		})
	}

	filter := bson.M{"_id": appointmentID}

	projection := bson.M{
		"_id":                        1,
		"serviceId":                  1,
		"pharmacy.information.name":  1,
		"pharmacy.information.image": 1,
		"pharmacy.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
		"facilityOrProfession":                                    1,
		"pharmacy.providerProvidedInformation.availableDrugs":     1,
		"pharmacy.providerProvidedInformation.notAvailableDrugs":  1,
		"pharmacy.providerProvidedInformation.homeDelivery":       1,
		"pharmacy.providerProvidedInformation.doctorApprovel":     1,
		"pharmacy.providerProvidedInformation.totalPriceToBePaid": 1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	var appointment entity.AppointmentEntity
	err = appointmentColl.FindOne(ctx, filter, findOptions).Decode(&appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pharmacy.GetPharmacyDrugsDetailResDto{
			Status:  false,
			Message: "Failed to fetch appointment data: " + err.Error(),
		})
	}

	var pharmacy1 entity.ServiceEntity
	reviewFilter := bson.M{"_id": appointment.ServiceID}
	projection = bson.M{
		"pharmacy.review.avgRating": 1,
	}

	reviewFindOptions := options.FindOne().SetProjection(projection)
	err = serviceColl.FindOne(ctx, reviewFilter, reviewFindOptions).Decode(&appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pharmacy.GetPharmacyDrugsDetailResDto{
			Status:  false,
			Message: "Failed to fetch average rating: " + err.Error(),
		})
	}

	var avgRating float64
	if pharmacy1.Pharmacy != nil {
		avgRating = pharmacy1.Pharmacy.Review.AvgRating
	}

	var pharmacyImage string
	var pharmacyName string
	var pharmacyAddress pharmacy.Address
	if appointment.Pharmacy != nil {
		pharmacyName = appointment.Pharmacy.Information.Name
		pharmacyImage = appointment.Pharmacy.Information.Image
		pharmacyAddress = pharmacy.Address(appointment.Pharmacy.Information.Address)
	}

	var availableDrugs string
	var notAvailableDrugs string
	var totalPriceToBePaid float64
	var homeDelivery string
	var doctorApprovel string
	if appointment.Pharmacy != nil && appointment.Pharmacy.ProvderProvidedInformation != nil {
		availableDrugs = appointment.Pharmacy.ProvderProvidedInformation.AvailableDrugs
		notAvailableDrugs = appointment.Pharmacy.ProvderProvidedInformation.NotAvailableDrugs
		totalPriceToBePaid = appointment.Pharmacy.ProvderProvidedInformation.TotalPriceToBePaid
		homeDelivery = appointment.Pharmacy.ProvderProvidedInformation.HomeDelivery
		doctorApprovel = appointment.Pharmacy.ProvderProvidedInformation.DoctorApprovel
	}

	expertiseRes := pharmacy.GetPharmacyDrugsDetailResDto{
		Status:  true,
		Message: "Data fetched successfully",
		Data: pharmacy.PharmacyDrugsRes{
			Id:                   appointment.Id,
			FacilityOrProfession: appointment.FacilityOrProfession,
			PharmacyInformation: pharmacy.PharmacyInformation{
				Id:        appointment.ServiceID,
				Name:      pharmacyName,
				Image:     pharmacyImage,
				Address:   pharmacyAddress,
				AvgRating: avgRating,
			},
			ProvderProvidedInformation: pharmacy.ProvderProvidedInformation{
				AvailableDrugs:     availableDrugs,
				NotAvailableDrugs:  notAvailableDrugs,
				HomeDelivery:       homeDelivery,
				TotalPriceToBePaid: totalPriceToBePaid,
				DoctorApprovel:     doctorApprovel,
			},
			// Customer: pharmacy.CustomerInformation{
			// 	Id:        appointment.Customer.ID,
			// 	FirstName: appointment.Customer.FirstName,
			// 	LastName:  appointment.Customer.LastName,
			// 	Image:     appointment.Customer.Image,
			// 	PhoneNumber: pharmacy.PhoneNumber{
			// 		DialCode:    appointment.Customer.PhoneNumber.DialCode,
			// 		Number:      appointment.Customer.PhoneNumber.Number,
			// 		CountryCode: appointment.Customer.PhoneNumber.CountryCode,
			// 	},
			// },
			// Prescription:         prescription,
			// NameAndQuantity:      nameAndQuantity,
			// ModeOfDelivery:       modeOfDelivery,
			// PricePaid:            pricePaid,
		},
	}

	return c.Status(fiber.StatusOK).JSON(expertiseRes)
}
