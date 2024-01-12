package providerAuth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	ProviderSignupReqDto struct {
		Email       string `json:"email" bson:"email"`
		DialCode    string `json:"dialCode" bson:"dialCode"`
		PhoneNumber string `json:"phoneNumber" bson:"phoneNumber"`
	}
	ProviderResponseDto struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}
)

type ProviderSignupVerifyOtpReqDto struct {
	Name        string `json:"name" bson:"name"`
	Email       string `json:"email" bson:"email"`
	DialCode    string `json:"dialCode" bson:"dialCode"`
	PhoneNumber string `json:"phoneNumber" bson:"phoneNumber"`
	CountryCode string `json:"countryCode" bson:"countryCode"`
	DeviceToken string `json:"deviceToken" bson:"deviceToken"`
	DeviceType  string `json:"deviceType" bson:"deviceType"`
	Password    string `json:"password" bson:"password"`
	EnteredOTP  string `json:"otp" bson:"otp"`
}

type ProviderSignupVerifyOtpResDto struct {
	Status   bool           `json:"status"`
	Message  string         `json:"message"`
	Token    string         `json:"token"`
	Provider ProviderResDto `json:"data"`
}

type ProviderResDto struct {
	Id                   primitive.ObjectID `json:"id" bson:"_id"`
	Name                 string             `json:"name" bson:"name"`
	Email                string             `json:"email" bson:"email"`
	Notification         Notification       `json:"notification" bson:"notification"`
	PhoneNumber          PhoneNumber        `json:"phoneNumber" bson:"phoneNumber"`
	Role                 string             `json:"role" bson:"role"`
	FacilityOrProfession string             `json:"facilityOrProfession" bson:"facilityOrProfession"`
	IsApproved           bool               `json:"isApproved" bson:"isApproved"`
	CreatedAt            time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt            time.Time          `json:"updatedAt" bson:"updatedAt"`
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
