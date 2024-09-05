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

// @Summary Get sended drugs by Id
// @Tags provider appointments
// @Description Get sended drugs by Id
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "appointment ID"
// @Produce json
// @Success 200 {object} services.GetPharmacyDrugsInfoResDto
// @Router /provider/services/appointment/get-sended-pharmacy-drug/{id} [get]
func GetSendedPharmacyDrugsByID(c *fiber.Ctx) error {
	
	var (
		appointmentColl = database.GetCollection("appointment")
	)

	idParam := c.Params("id")
	appointmentID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.GetPharmacyDrugsInfoResDto{
			Status:  false,
			Message: "Invalid appointment ID",
		})
	}

	filter := bson.M{"_id": appointmentID}

	projection := bson.M{
		"_id": 1,
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
		return c.Status(fiber.StatusInternalServerError).JSON(services.GetPharmacyDrugsInfoResDto{
			Status:  false,
			Message: "Failed to fetch appointment data: " + err.Error(),
		})
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

	expertiseRes := services.GetPharmacyDrugsInfoResDto{
		Status:  true,
		Message: "Data fetched successfully",
		Data: services.PharmacyDrugsInfoRes{
			AppointmentId:      appointment.Id,
			AvailableDrugs:     availableDrugs,
			NotAvailableDrugs:  notAvailableDrugs,
			HomeDelivery:       homeDelivery,
			DoctorApprovel:     doctorApprovel,
			TotalPriceToBePaid: totalPriceToBePaid,
		},
	}

	return c.Status(fiber.StatusOK).JSON(expertiseRes)
}
