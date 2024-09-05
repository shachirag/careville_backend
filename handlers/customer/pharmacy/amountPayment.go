package pharmacy

import (
	"time"

	"careville_backend/database"
	"careville_backend/dto/customer/pharmacy"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// @Summary amount payment for pharmacy
// @Tags customer pharmacy
// @Description amount payment for pharmacy
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param customer body pharmacy.AmountPaymentForPharmacyDrugsReqDto true "send pharmacy drugs"
// @Produce json
// @Success 200 {object} pharmacy.AmountPaymentForPharmacyDrugsResDto
// @Router /customer/healthFacility/appointment/pharmacy-amount-payment [put]
func AmountPaymentForPharmacy(c *fiber.Ctx) error {
	var (
		appointmentColl = database.GetCollection("appointment")
		data            pharmacy.AmountPaymentForPharmacyDrugsReqDto
	)

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pharmacy.AmountPaymentForPharmacyDrugsResDto{
			Status:  false,
			Message: "Failed to parse request body: " + err.Error(),
		})
	}

	pharmacyUpdate := bson.M{
		"$set": bson.M{
			"pharmacy.pricePaid": data.Amount,
			"updatedAt":          time.Now().UTC(),
		},
	}

	filter := bson.M{"_id": data.AppointmentId}

	_, err = appointmentColl.UpdateOne(ctx, filter, pharmacyUpdate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pharmacy.AmountPaymentForPharmacyDrugsResDto{
			Status:  false,
			Message: "Failed to paid ammount into MongoDB: " + err.Error(),
		})
	}

	pharmacyRes := pharmacy.AmountPaymentForPharmacyDrugsResDto{
		Status:  true,
		Message: "Successfully paid amount",
	}
	return c.Status(fiber.StatusOK).JSON(pharmacyRes)
}
