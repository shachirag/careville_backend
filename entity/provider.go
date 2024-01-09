package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	ProviderEntity struct {
		Id                   primitive.ObjectID `json:"_id" bson:"_id"`
		Name                 string             `json:"name" bson:"name"`
		Email                string             `json:"email" bson:"email"`
		Password             string             `json:"password" bson:"password"`
		Image                string             `json:"image" bson:"image"`
		CreatedAt            time.Time          `json:"createdAt" bson:"createdAt"`
		UpdatedAt            time.Time          `json:"updatedAt" bson:"updatedAt"`
		PhoneNumber          PhoneNumber        `json:"phoneNumber" bson:"phoneNumber"`
		AdditionalDetails    string             `json:"additionalDetails" bson:"additionalDetails"`
		Address              string             `json:"address" bson:"address"`
		Notification         bool               `json:"notification" bson:"notification"`
		IsEmergencyAvailable bool               `json:"isEmergencyAvailable" bson:"isEmergencyAvailable"`
	}
	PhoneNumber struct {
		DialCode string `json:"dialCode" bson:"dialCode"`
		Number   string `json:"number" bson:"number"`
	}
)
