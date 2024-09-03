package commonApi

import (
	"careville_backend/database"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"
	"careville_backend/utils"
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
// @Param reason query string false "reason for rejection"
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
	reason := c.Query("reason")

	appointmentColl := database.GetCollection("appointment")

	filter := bson.M{"_id": appointmentObjID}
	var appointment entity.AppointmentEntity
	err = appointmentColl.FindOne(ctx, filter).Decode(&appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.NotificationResDto{
			Status:  false,
			Message: "Failed to fetch appointment data: " + err.Error(),
		})
	}

	var customer entity.CustomerEntity
	customerFilter := bson.M{
		"_id": appointment.Customer.ID,
	}
	err = database.GetCollection("customer").FindOne(ctx, customerFilter).Decode(&customer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.NotificationResDto{
			Status:  false,
			Message: "Failed to fetch customer data: " + err.Error(),
		})
	}

	var update bson.M
	var notificationTitle, notificationBody string

	if isEnableParam == "approved" {
		update = bson.M{"$set": bson.M{"appointmentStatus": "approved"}}
		notificationTitle = "Appointment Confirmed"
		notificationBody = "Your appointment has been confirmed."
	} else if isEnableParam == "rejected" {
		update = bson.M{"$set": bson.M{"appointmentStatus": "rejected"}}
		notificationTitle = "Appointment Rejected"
		notificationBody = "Your appointment has been rejected. Reason: " + reason
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(services.NotificationResDto{
			Status:  false,
			Message: "Invalid status",
		})
	}

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

	if customer.Notification.DeviceToken != "" && customer.Notification.DeviceType != "" {
		notificationData := map[string]string{
			"type":                 "booking-update",
			"appointmentId":        appointmentId,
			"role":                 appointment.Role,
			"facilityOrProfession": appointment.FacilityOrProfession,
		}

		utils.SendNotificationToUser(customer.Notification.DeviceToken, customer.Notification.DeviceType, notificationTitle, notificationBody, notificationData, customer.Id, "customer")
	}

	response := services.NotificationResDto{
		Status:  true,
		Message: "success",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
