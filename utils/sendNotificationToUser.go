package utils

import (
	"careville_backend/database"
	"careville_backend/entity"
	"careville_backend/firebase"
	"context"

	"firebase.google.com/go/messaging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ctx = context.Background()

func SendNotificationToUser(
	deviceToken string,
	deviceType string,
	title string,
	body string,
	data map[string]string,
	userId primitive.ObjectID,
	role string,
) {

	var notificationEnabled bool
	if role == "customer" {
		var customer entity.CustomerEntity
		customerFilter := bson.M{"_id": userId}
		database.GetCollection("customer").FindOne(ctx, customerFilter).Decode(&customer)
		notificationEnabled = customer.Notification.IsEnabled
	} else if role == "provider" {
		var service entity.ServiceEntity
		serviceFilter := bson.M{"_id": userId}
		database.GetCollection("service").FindOne(ctx, serviceFilter).Decode(&service)
		notificationEnabled = service.User.Notification.IsEnabled
	}

	if notificationEnabled {
		var message *messaging.Message = &messaging.Message{
			Data: data,
			Notification: &messaging.Notification{
				Title: title,
				Body:  body,
			},
			Token: deviceToken,
		}
		firebase.GetFirebaseMessagingClient().Send(ctx, message)
	}
	SaveNotification(userId, title, body, data)
}
