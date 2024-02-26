package fitnessCenter

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetFitnessCenterAppointmentDetailResDto struct {
	Status  bool                        `json:"status"`
	Message string                      `json:"message"`
	Data    FitnessCenterAppointmentRes `json:"data"`
}

type FitnessCenterAppointmentRes struct {
	Id                       primitive.ObjectID       `json:"id" bson:"id"`
	FitnessCenterInformation FitnessCenterInformation `json:"fitnessCenterInformation" bson:"fitnessCenterInformation"`
	Customer                 CustomerInformation      `json:"customer" bson:"customer"`
	TrainerInformation       TrainerInformation       `json:"trainer" bson:"trainer"`
	Subscription             SubscriptionData         `json:"subscriptionType" bson:"subscriptionType"`
	FacilityOrProfession     string                   `json:"facilityOrProfession" bson:"facilityOrProfession"`
	PricePaid                float64                  `json:"pricePaid" bson:"pricePaid"`
	FamilyMember             FamilyMember             `json:"familyMember"`
}

type FitnessCenterInformation struct {
	Id        primitive.ObjectID `json:"id" bson:"id"`
	Name      string             `json:"name" bson:"name"`
	Address   Address            `json:"address" bson:"address"`
	Image     string             `json:"image" bson:"image"`
	AvgRating float64            `json:"avgRating" bson:"avgRating"`
}

type SubscriptionData struct {
	Package string  `json:"package"`
	Price   float64 `json:"price"`
}

type FamilyMember struct {
	Id           primitive.ObjectID `json:"id" bson:"id"`
	Name         string             `json:"name" bson:"name"`
	Age          string             `json:"age" bson:"age"`
	Sex          string             `json:"sex" bson:"sex"`
	RelationShip string             `json:"relationship" bson:"relationship"`
}

type CustomerInformation struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	FirstName   string             `json:"firstName" bson:"firstName"`
	LastName    string             `json:"lastName" bson:"lastName"`
	Image       string             `json:"image" bson:"image"`
	PhoneNumber PhoneNumber        `json:"phoneNumber" bson:"phoneNumber"`
}

type PhoneNumber struct {
	DialCode    string `json:"dialCode" bson:"dialCode"`
	Number      string `json:"number" bson:"number"`
	CountryCode string `json:"countryCode" bson:"countryCode"`
}

type TrainerInformation struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Category    string             `json:"category" bson:"category"`
	Name        string             `json:"name" bson:"name"`
	Information string             `json:"information" bson:"information"`
	Price       float64            `json:"price" bson:"price"`
}
