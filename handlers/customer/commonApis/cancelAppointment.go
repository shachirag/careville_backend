package common

import (
	"careville_backend/database"
	common "careville_backend/dto/customer/commonApis"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// @Summary cancel appointments
// @Tags customer commonApis
// @Description cancel appointments
// @Accept json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param appointmentId query string true "appointment Id"
// @Produce json
// @Success 200 {object} common.CancelAppointmentResDto
// @Router /customer/cancel-appointment [put]
func CancelAppointment(c *fiber.Ctx) error {

	appointmentId := c.Query("appointmentId")
	appointmentObjID, err := primitive.ObjectIDFromHex(appointmentId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.CancelAppointmentResDto{
			Status:  false,
			Message: "Invalid appointment ID",
		})
	}

	appointmentColl := database.GetCollection("appointment")

	filter := bson.M{"_id": appointmentObjID}
	var appointment entity.AppointmentEntity
	err = appointmentColl.FindOne(ctx, filter).Decode(&appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.CancelAppointmentResDto{
			Status:  false,
			Message: "Failed to fetch appointment data: " + err.Error(),
		})
	}

	update := bson.M{"$set": bson.M{"appointmentStatus": "cancelled"}}

	result, err := appointmentColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.CancelAppointmentResDto{
			Status:  false,
			Message: "Failed to cancel appointment: " + err.Error(),
		})
	}

	if result.ModifiedCount == 0 {
		return c.Status(fiber.StatusOK).JSON(common.CancelAppointmentResDto{
			Status:  false,
			Message: "No documents were modified",
		})
	}

	response := common.CancelAppointmentResDto{
		Status:  true,
		Message: "success",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
