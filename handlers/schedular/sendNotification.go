package schedular

import (
	"careville_backend/database"
	"careville_backend/entity"
	"careville_backend/utils"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SendAppointmentReminderNotification(c *fiber.Ctx) error {
	var (
		appColl = database.GetCollection("appointment")
		ctx     = context.Background()
	)

	currentTime := time.Now().UTC()

	startTime := currentTime.Add(15 * time.Minute)
	endTime := startTime.Add(15 * time.Minute)

	filter := bson.M{
		"$or": []bson.M{
			{
				"physiotherapist.appointmentDetails.from": bson.M{
					"$gt":  startTime,
					"$lte": endTime,
				},
				"appointmentStatus": "approved",
			},
			{
				"hospital.appointmentDetails.from": bson.M{
					"$gt":  startTime,
					"$lte": endTime,
				},
				"appointmentStatus": "approved",
			},
			{
				"laboratory.appointmentDetails.date": bson.M{
					"$gt":  startTime,
					"$lte": endTime,
				},
				"appointmentStatus": "approved",
			},
			{
				"nurse.appointmentDetails.from": bson.M{
					"$gt":  startTime,
					"$lte": endTime,
				},
				"appointmentStatus": "approved",
			},
			{
				"medicalLabScientist.appointmentDetails.from": bson.M{
					"$gt":  startTime,
					"$lte": endTime,
				},
				"appointmentStatus": "approved",
			},
			{
				"doctor.appointmentDetails.from": bson.M{
					"$gt":  startTime,
					"$lte": endTime,
				},
				"appointmentStatus": "approved",
			},
		},
	}

	cursor, err := appColl.Find(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(NotificationRes{
			Status:  false,
			Message: "Failed to retrieve appointments: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var appointments []entity.AppointmentEntity
	var customerIds []primitive.ObjectID
	var providerIds []primitive.ObjectID

	for cursor.Next(ctx) {
		var appointment entity.AppointmentEntity
		if err := cursor.Decode(&appointment); err != nil {
			continue
		}

		appointments = append(appointments, appointment)

		if appointment.Customer.ID != primitive.NilObjectID {
			customerIds = append(customerIds, appointment.Customer.ID)
		}

		if appointment.ServiceID != primitive.NilObjectID {
			providerIds = append(providerIds, appointment.ServiceID)
		}

	}

	if err := cursor.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(NotificationRes{
			Status:  false,
			Message: "Error iterating over appointments: " + err.Error(),
		})
	}

	var customerDetails map[primitive.ObjectID]entity.CustomerEntity
	if len(customerIds) > 0 {
		customerDetails, err = FetchCustomerDetails(ctx, customerIds)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(NotificationRes{
				Status:  false,
				Message: "Failed to retrieve customer details: " + err.Error(),
			})
		}
	}

	var providerDetails map[primitive.ObjectID]entity.ServiceEntity
	if len(providerIds) > 0 {
		providerDetails, err = FetchProviderDetails(ctx, providerIds)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(NotificationRes{
				Status:  false,
				Message: "Failed to retrieve provider details: " + err.Error(),
			})
		}
	}

	var wg sync.WaitGroup

	for _, appointment := range appointments {
		customerData, okCus := customerDetails[appointment.Customer.ID]
		providerData, okPro := providerDetails[appointment.ServiceID]
		if !okPro || !okCus {
			continue
		}

		appointmentDetails := "Your appointment will start in 15 minutes."

		wg.Add(1)
		go func(appointment entity.AppointmentEntity, customerData entity.CustomerEntity, providerData entity.ServiceEntity) {
			defer wg.Done()

			title := "Appointment Alarm"
			body := fmt.Sprintf("%v", appointmentDetails)
			data := map[string]string{
				"type":                 "Appointment Reminder",
				"appointmentId":        appointment.Id.Hex(),
				"facilityOrProfession": appointment.FacilityOrProfession,
				"role":                 appointment.Role,
			}

			if customerData.Notification.DeviceToken != "" && customerData.Notification.DeviceType != "" {
				utils.SendNotificationToUser(
					customerData.Notification.DeviceToken,
					customerData.Notification.DeviceType,
					title,
					body,
					data,
					appointment.Customer.ID,
					"customer",
				)
			}

			if providerData.User.Notification.DeviceToken != "" && providerData.User.Notification.DeviceType != "" {
				utils.SendNotificationToUser(
					providerData.User.Notification.DeviceToken,
					providerData.User.Notification.DeviceType,
					title,
					body,
					data,
					appointment.ServiceID,
					"provider",
				)
			}
		}(appointment, customerData, providerData)
	}

	wg.Wait()

	return c.Status(fiber.StatusOK).JSON(NotificationRes{
		Status:  true,
		Message: "Notifications sent.",
	})
}

func FetchCustomerDetails(ctx context.Context, customerIDs []primitive.ObjectID) (map[primitive.ObjectID]entity.CustomerEntity, error) {
	customerDetails := make(map[primitive.ObjectID]entity.CustomerEntity)
	if len(customerIDs) > 0 {
		customerFilter := bson.M{"_id": bson.M{"$in": customerIDs}}
		customerCur, err := database.GetCollection("customer").Find(ctx, customerFilter)
		if err != nil {
			return nil, err
		}
		defer customerCur.Close(ctx)

		for customerCur.Next(ctx) {
			var customer entity.CustomerEntity
			err := customerCur.Decode(&customer)
			if err != nil {
				return nil, err
			}
			customerDetails[customer.Id] = customer
		}
	}

	return customerDetails, nil
}

func FetchProviderDetails(ctx context.Context, providerIDs []primitive.ObjectID) (map[primitive.ObjectID]entity.ServiceEntity, error) {
	providerDetails := make(map[primitive.ObjectID]entity.ServiceEntity)
	if len(providerIDs) > 0 {
		providerFilter := bson.M{"_id": bson.M{"$in": providerIDs}}
		providerCur, err := database.GetCollection("service").Find(ctx, providerFilter)
		if err != nil {
			return nil, err
		}
		defer providerCur.Close(ctx)

		for providerCur.Next(ctx) {
			var provider entity.ServiceEntity
			err := providerCur.Decode(&provider)
			if err != nil {
				return nil, err
			}
			providerDetails[provider.Id] = provider
		}
	}

	return providerDetails, nil
}
