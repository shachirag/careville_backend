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
	Id                primitive.ObjectID `json:"id" bson:"_id"`
	Name              string             `json:"name" bson:"name"`
	Email             string             `json:"email" bson:"email"`
	Image             string             `json:"image" bson:"image"`
	CreatedAt         time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt         time.Time          `json:"updatedAt" bson:"updatedAt"`
	PhoneNumber       PhoneNumber        `json:"phoneNumber" bson:"phoneNumber"`
	AdditionalDetails string             `json:"additionalDetails" bson:"additionalDetails"`
	Address           string             `json:"address" bson:"address"`
	Notification      bool               `json:"notification" bson:"notification"`
}
