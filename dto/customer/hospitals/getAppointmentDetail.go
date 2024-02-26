package hospitals

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetHospitalAppointmentDetailResDto struct {
	Status  bool                   `json:"status"`
	Message string                 `json:"message"`
	Data    HospitalAppointmentRes `json:"data"`
}

type HospitalAppointmentRes struct {
	Id                   primitive.ObjectID  `json:"id" bson:"id"`
	HospitalInformation  HospitalInformation `json:"hospitalInformation" bson:"hospitalInformation"`
	Customer             CustomerInformation `json:"customer" bson:"customer"`
	AppointmentDetails   AppointmentDetails  `json:"appointmentDetails" bson:"appointmentDetails"`
	FacilityOrProfession string              `json:"facilityOrProfession" bson:"facilityOrProfession"`
	PricePaid            float64             `json:"pricePaid" bson:"pricePaid"`
	FamilyMember         FamilyMember        `json:"familyMember"`
}

type HospitalInformation struct {
	Id        primitive.ObjectID `json:"id" bson:"id"`
	Name      string             `json:"name" bson:"name"`
	Address   Address            `json:"address" bson:"address"`
	Image     string             `json:"image" bson:"image"`
	AvgRating float64            `json:"avgRating" bson:"avgRating"`
}

type FamilyMember struct {
	Id           primitive.ObjectID `json:"id" bson:"id"`
	Name         string             `json:"name" bson:"name"`
	Age          string             `json:"age" bson:"age"`
	Sex          string             `json:"sex" bson:"sex"`
	RelationShip string             `json:"relationship" bson:"relationship"`
}

type AppointmentDetails struct {
	AppointmentFromDate time.Time `json:"appointmentFromDate" bson:"appointmentFromDate"`
	AppointmentToDate   time.Time `json:"appointmentToDate" bson:"appointmentToDate"`
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
	CountryCode string `json:"countryCode" bson:"countryCode"`
	Number      string `json:"number" bson:"number"`
}
