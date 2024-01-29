package adminAuth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginAdminReqDto struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type LoginAdminResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type LoginVerifyOtpReqDto struct {
	Email      string `json:"email" bson:"email"`
	EnteredOTP string `json:"otp" bson:"otp"`
}

type LoginVerifyOtpResDto struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    GetAdminRes `json:"data"`
	Token   string      `json:"token"`
}

type GetAdminRes struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	FirstName string             `json:"firstName" bson:"firstName"`
	LastName  string             `json:"lastName" bson:"lastName"`
	Email     string             `json:"email" bson:"email"`
	Image     string             `json:"image" bson:"image"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}
