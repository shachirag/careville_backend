package providerAuth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetProviderResDto struct {
	Status   bool         `json:"status"`
	Message  string       `json:"message"`
	Provider ProviderData `json:"data"`
}

type ProviderData struct {
	Id                   primitive.ObjectID `json:"id" bson:"_id"`
	FirstName            string             `json:"firstName" bson:"firstName"`
	LastName             string             `json:"lastName" bson:"lastName"`
	Email                string             `json:"email" bson:"email"`
	Image                string             `json:"image" bson:"image"`
	CreatedAt            time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt            time.Time          `json:"updatedAt" bson:"updatedAt"`
	PhoneNumber          PhoneNumber        `json:"phoneNumber" bson:"phoneNumber"`
	AdditionalDetails    string             `json:"additionalDetails" bson:"additionalDetails"`
	Address              Address            `json:"address" bson:"address"`
	IsEmergencyAvailable bool               `json:"isEmergencyAvailable" bson:"isEmergencyAvailable"`
	Notification         Notification       `json:"notification" bson:"notification"`
	Documents            Documents          `json:"documents" bson:"documents"`
}

type Documents struct {
	Certificate string `json:"certificate" bson:"certificate"`
	License     string `json:"license" bson:"license"`
}
type Address struct {
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
	Add         string    `json:"add" bson:"add"`
	Type        string    `json:"type" bson:"type"`
}
