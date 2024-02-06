package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewEntity struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Customer    Customer           `json:"customer" bson:"customer"`
	ServiceId   primitive.ObjectID `json:"serviceId" bson:"serviceId"`
	Description string           `json:"description" bson:"description"`
	Rating      float64            `json:"rating" bson:"rating"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type Customer struct {
	Id        primitive.ObjectID `json:"id" bson:"id"`
	FirstName string             `json:"firstName" bson:"firstName"`
	LastName  string             `json:"lastName" bson:"lastName"`
	Image     string             `json:"image" bson:"image"`
}
