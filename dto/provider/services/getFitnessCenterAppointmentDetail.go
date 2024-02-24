package services

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetFitnessCenterAppointmentDetailResDto struct {
	Status  bool                   `json:"status"`
	Message string                 `json:"message"`
	Data    FitnessCenterAppointmentRes `json:"data"`
}

type FitnessCenterAppointmentRes struct {
	Id                   primitive.ObjectID  `json:"id" bson:"id"`
	Customer             CustomerInformation `json:"customer" bson:"customer"`
	TrainerInformation   TrainerInformation  `json:"trainer" bson:"trainer"`
	Subscription         SubscriptionData    `json:"subscriptionType" bson:"subscriptionType"`
	FacilityOrProfession string              `json:"facilityOrProfession" bson:"facilityOrProfession"`
	PricePaid            float64             `json:"pricePaid" bson:"pricePaid"`
	FamilyMember         FamilyMember        `json:"familyMember"`
}

type SubscriptionData struct {
	Package string  `json:"package"`
	Price   float64 `json:"price"`
}

type TrainerInformation struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Category    string             `json:"category" bson:"category"`
	Name        string             `json:"name" bson:"name"`
	Information string             `json:"information" bson:"information"`
	Price       float64            `json:"price" bson:"price"`
}
