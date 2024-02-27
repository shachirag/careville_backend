package pharmacy

import (
	"careville_backend/database"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get appointment by ID
// @Tags provider appointments
// @Description Get appointment by ID
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "appointment ID"
// @Produce json
// @Success 200 {object} services.GetPharmacyDrugsDetailResDto
// @Router /provider/services/appointment/pharmacy-drug/{id} [get]
func GetPharmacyAppointmentByID(c *fiber.Ctx) error {
	var (
		appointmentColl = database.GetCollection("appointment")
	)

	idParam := c.Params("id")
	appointmentID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.GetPharmacyDrugsDetailResDto{
			Status:  false,
			Message: "Invalid appointment ID",
		})
	}

	filter := bson.M{"_id": appointmentID}

	projection := bson.M{
		"_id":                1,
		"customer.id":        1,
		"customer.firstName": 1,
		"customer.lastName":  1,
		"customer.image":     1,
		"customer.phoneNumber": bson.M{
			"dialCode":    1,
			"number":      1,
			"countryCode": 1,
		},
		"facilityOrProfession":                    1,
		"pharmacy.pricePaid":                      1,
		"pharmacy.requestedDrugs.nameAndQuantity": 1,
		"pharmacy.requestedDrugs.modeOfDelivery": 1,
		"pharmacy.requestedDrugs.prescription":    1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	var appointment entity.AppointmentEntity
	err = appointmentColl.FindOne(ctx, filter, findOptions).Decode(&appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.GetPharmacyDrugsDetailResDto{
			Status:  false,
			Message: "Failed to fetch appointment data: " + err.Error(),
		})
	}

	var nameAndQuantity string
	var modeOfDelivery string
	var pricePaid float64
	var prescription []string
	if appointment.Pharmacy != nil {
		nameAndQuantity = appointment.Pharmacy.RequestedDrugs.NameAndQuantity
		modeOfDelivery = appointment.Pharmacy.RequestedDrugs.ModeOfDelivery
		pricePaid = appointment.Pharmacy.PricePaid
		prescription = appointment.Pharmacy.RequestedDrugs.Prescription
	}

	expertiseRes := services.GetPharmacyDrugsDetailResDto{
		Status:  true,
		Message: "Data fetched successfully",
		Data: services.PharmacyDrugsRes{
			Id: appointment.Id,
			Customer: services.CustomerInformation{
				Id:        appointment.Customer.ID,
				FirstName: appointment.Customer.FirstName,
				LastName:  appointment.Customer.LastName,
				Image:     appointment.Customer.Image,
				PhoneNumber: services.PhoneNumber{
					DialCode:    appointment.Customer.PhoneNumber.DialCode,
					Number:      appointment.Customer.PhoneNumber.Number,
					CountryCode: appointment.Customer.PhoneNumber.CountryCode,
				},
			},
			Prescription:         prescription,
			NameAndQuantity:      nameAndQuantity,
			ModeOfDelivery:       modeOfDelivery,
			FacilityOrProfession: appointment.FacilityOrProfession,
			PricePaid:            pricePaid,
		},
	}

	return c.Status(fiber.StatusOK).JSON(expertiseRes)
}
