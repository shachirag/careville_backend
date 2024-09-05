package pharmacy

import (
	"encoding/json"
	"time"

	"careville_backend/database"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// @Summary send pharmcy frugs
// @Tags provider appointments
// @Description Summary send pharmcy frugs
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param provider body services.SendPharmacyDrugsInfoReqDto true "send pharmacy drugs"
// @Produce json
// @Success 200 {object} services.SendPharmacyDrugsInfoResDto
// @Router /provider/services/appointment/send-pharmacy-drugs [put]
func SendPharmacyDrugs(c *fiber.Ctx) error {
	var (
		appointmentColl = database.GetCollection("appointment")
		data            services.SendPharmacyDrugsInfoReqDto
	)

	dataStr := c.FormValue("data")
	dataBytes := []byte(dataStr)

	err := json.Unmarshal(dataBytes, &data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.SendPharmacyDrugsInfoResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	pharmacyDrugsData := entity.ProvderProvidedInformation{
		AvailableDrugs:     data.AvailableDrugs,
		NotAvailableDrugs:  data.NotAvailableDrugs,
		TotalPriceToBePaid: data.TotalPriceToBePaid,
		HomeDelivery:       data.HomeDelivery,
		DoctorApprovel:     data.DoctorApprovel,
	}

	pharmacyUpdate := bson.M{
		"$set": bson.M{
			"pharmacy.provderProvidedInformation": pharmacyDrugsData,
			"updatedAt":                           time.Now().UTC(),
		},
	}

	filter := bson.M{"_id": data.AppointmentId}

	_, err = appointmentColl.UpdateOne(ctx, filter, pharmacyUpdate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.SendPharmacyDrugsInfoResDto{
			Status:  false,
			Message: "Failed to send drugs data into MongoDB: " + err.Error(),
		})
	}

	pharmacyRes := services.SendPharmacyDrugsInfoResDto{
		Status:  true,
		Message: "Successfully sent drugs",
	}
	return c.Status(fiber.StatusOK).JSON(pharmacyRes)
}
