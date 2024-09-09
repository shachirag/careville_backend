package utils

import (
	"careville_backend/database"
	"careville_backend/entity"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SaveNotification(
	userId primitive.ObjectID,
	title string,
	body string,
	data map[string]string,
) error {
	var notification = entity.NotificationEntity{
		Id:        primitive.NewObjectID(),
		UserId:    userId,
		Title:     title,
		Body:      body,
		Data:      data,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	_, err := database.GetCollection("notification").InsertOne(ctx, notification)
	if err != nil {
		return fmt.Errorf("error inserting notification: %v", err)
	}

	return nil
}
