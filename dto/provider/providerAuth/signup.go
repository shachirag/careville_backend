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
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Email       string             `json:"email" bson:"email"`
	Image       string             `json:"image" bson:"image"`
	PhoneNumber PhoneNumber        `json:"phoneNumber" bson:"phoneNumber"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type PhoneNumber struct {
	DialCode string `json:"dialCode" bson:"dialCode"`
	Number   string `json:"number" bson:"number"`
}
