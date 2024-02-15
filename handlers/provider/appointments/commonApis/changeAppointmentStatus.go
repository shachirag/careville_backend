package commonApi

import (
	"careville_backend/database"
	"careville_backend/dto/provider/services"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ctx = context.Background()

// @Summary change appointment status
// @Tags provider appointments
// @Description change appointment status
// @Accept json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param appointmentId query string true "appointment Id"
// @Param status query string false "change status"
// @Produce json
// @Success 200 {object} services.NotificationResDto
// @Router /provider/appointment/change-appointment-status [put]
func ChangeAppointmentStatus(c *fiber.Ctx) error {

	appointmentId := c.Query("appointmentId")
	appointmentObjID, err := primitive.ObjectIDFromHex(appointmentId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.NotificationResDto{
			Status:  false,
			Message: "Invalid appointment ID",
		})
	}
	
	isEnableParam := c.Query("status")

	appointmentColl := database.GetCollection("appointment")

	filter := bson.M{"_id": appointmentObjID}

	update := bson.M{"$set": bson.M{"appointmentStatus": isEnableParam}}

	result, err := appointmentColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.NotificationResDto{
			Status:  false,
			Message: "Failed to change status: " + err.Error(),
		})
	}

	if result.ModifiedCount == 0 {
		return c.Status(fiber.StatusOK).JSON(services.NotificationResDto{
			Status:  false,
			Message: "No documents were modified",
		})
	}

	response := services.NotificationResDto{
		Status:  true,
		Message: "success",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
