package providerAuth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginProviderReqDto struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type LoginProviderResDto struct {
	Status   bool            `json:"status"`
	Message  string          `json:"message"`
	Provider ProviderRespDto `json:"data"`
	Token    string          `json:"token"`
}

type ProviderRespDto struct {
	Role Role `json:"role"`
	User User `json:"user"`
}

type Role struct {
	ProviderId           primitive.ObjectID `json:"providerId" bson:"providerId"`
	Role                 string             `json:"role" bson:"role"`
	FacilityOrProfession string             `json:"facilityOrProfession" bson:"facilityOrProfession"`
	Status               string             `json:"status" bson:"status"`
}

type User struct {
	Id           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Email        string             `json:"email" bson:"email"`
	Image        string             `json:"image" bson:"image"`
	Notification Notification       `json:"notification" bson:"notification"`
	PhoneNumber  PhoneNumber        `json:"phoneNumber" bson:"phoneNumber"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type GetProvideResDto struct {
	Status   bool           `json:"status"`
	Message  string         `json:"message"`
	Provider ProviderResDto `json:"data"`
}
