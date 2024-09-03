package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationEntity struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	UserId    primitive.ObjectID `json:"customerId" bson:"customerId"`
	Title     string             `json:"title" bson:"title"`
	Body      string             `json:"body" bson:"body"`
	Data      map[string]string  `json:"data" bson:"data"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
