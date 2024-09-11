package laboratory

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetLaboratoryAppointmentDetailResDto struct {
	Status  bool                     `json:"status"`
	Message string                   `json:"message"`
	Data    LaboratoryAppointmentRes `json:"data"`
}

type LaboratoryAppointmentRes struct {
	LaboratoryInformation LaboratoryInformation `json:"laboratoryInformation" bson:"laboratoryInformation"`
	Id                    primitive.ObjectID    `json:"id" bson:"id"`
	Customer              CustomerInformation   `json:"customer" bson:"customer"`
	AppointmentDetails    AppointmentData       `json:"appointmentDetails" bson:"appointmentDetails"`
	Investigation         Investigation         `json:"investigation" bson:"investigation"`
	FacilityOrProfession  string                `json:"facilityOrProfession" bson:"facilityOrProfession"`
	PricePaid             float64               `json:"pricePaid" bson:"pricePaid"`
	FamilyMember          FamilyMember          `json:"familyMember"`
}

type LaboratoryInformation struct {
	Id        primitive.ObjectID `json:"id" bson:"id"`
	Name      string             `json:"name" bson:"name"`
	Address   Address            `json:"address" bson:"address"`
	Image     string             `json:"image" bson:"image"`
	AvgRating float64            `json:"avgRating" bson:"avgRating"`
}

type AppointmentData struct {
	AppointmentDate time.Time `json:"appointmentDate" bson:"appointmentDate"`
}

type Investigation struct {
	ID          primitive.ObjectID `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	Information string             `json:"information" bson:"information"`
	Type        string             `json:"type" bson:"type"`
	Price       float64            `json:"price" bson:"price"`
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
	Age         string             `json:"age" bson:"age"`
}

type PhoneNumber struct {
	DialCode    string `json:"dialCode" bson:"dialCode"`
	Number      string `json:"number" bson:"number"`
	CountryCode string `json:"countryCode" bson:"countryCode"`
}
