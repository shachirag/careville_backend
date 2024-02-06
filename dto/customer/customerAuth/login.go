package customerAuth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginCustomerReqDto struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type LoginCustomerResDto struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Data    CustomerResDto `json:"data"`
	Token   string         `json:"token"`
}

type GetCustomerResDto struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Data    CustomerResDto `json:"data"`
}

type CustomerResDto struct {
	Id            primitive.ObjectID `json:"id" bson:"_id"`
	FirstName     string             `json:"firstName" bson:"firstName"`
	LastName      string             `json:"lastName" bson:"lastName"`
	Email         string             `json:"email" bson:"email"`
	Image         string             `json:"image" bson:"image"`
	PhoneNumber   PhoneNumber        `json:"phoneNumber" bson:"phoneNumber"`
	Notification  Notification       `json:"notification" bson:"notification"`
	FamilyMembers []FamilyMembers    `json:"familyMembers" bson:"familyMembers"`
	Sex           string             `json:"sex" bson:"sex"`
	Age           string             `json:"age" bson:"age"`
	Address       Address            `json:"address" bson:"address"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type FamilyMembers struct {
	Id           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Age          string             `json:"age" bson:"age"`
	Sex          string             `json:"sex" bson:"sex"`
	RelationShip string             `json:"relationShip" bson:"relationShip"`
}

type Address struct {
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
	Add         string    `json:"add" bson:"add"`
	Type        string    `json:"type" bson:"type"`
}

type Notification struct {
	DeviceToken string `json:"deviceToken" bson:"deviceToken"`
	DeviceType  string `json:"deviceType" bson:"deviceType"`
	IsEnabled   bool   `json:"isEnabled" bson:"isEnabled"`
}

type PhoneNumber struct {
	DialCode    string `json:"dialCode" bson:"dialCode"`
	CountryCode string `json:"countryCode" bson:"countryCode"`
	Number      string `json:"number" bson:"number"`
}
