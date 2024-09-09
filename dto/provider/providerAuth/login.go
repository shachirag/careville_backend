package providerAuth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginProviderReqDto struct {
	Email       string `json:"email" bson:"email"`
	Password    string `json:"password" bson:"password"`
	DeviceToken string `json:"deviceToken" bson:"deviceToken"`
	DeviceType  string `json:"deviceType" bson:"deviceType"`
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
	Role                 string `json:"role" bson:"role"`
	FacilityOrProfession string `json:"facilityOrProfession" bson:"facilityOrProfession"`
	ServiceStatus        string `json:"serviceStatus" bson:"serviceStatus"`
	Name                 string `json:"name" bson:"name"`
	Image                string `json:"image" bson:"image"`
	IsEmergencyAvailable bool   `json:"isEmergencyAvailable" bson:"isEmergencyAvailable"`
}

type User struct {
	Id           primitive.ObjectID `json:"id" bson:"_id"`
	FirstName    string             `json:"firstName" bson:"firstName"`
	LastName     string             `json:"lastName" bson:"lastName"`
	Email        string             `json:"email" bson:"email"`
	Notification Notification       `json:"notification" bson:"notification"`
	PhoneNumber  PhoneNumber        `json:"phoneNumber" bson:"phoneNumber"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type GetProvideResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	// Provider ProviderResDto `json:"data"`
}
