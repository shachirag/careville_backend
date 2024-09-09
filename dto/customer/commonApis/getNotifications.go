package common

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationResData struct {
	Status  bool                 `json:"status"`
	Message string               `json:"message"`
	Data    []GetNotificationRes `json:"data"`
}

type GetNotificationRes struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Title     string             `json:"title" bson:"title"`
	Body      string             `json:"body" bson:"body"`
	Data      map[string]string  `json:"data" bson:"data"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
