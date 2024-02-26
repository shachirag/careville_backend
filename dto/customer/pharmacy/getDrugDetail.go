package pharmacy

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetPharmacyDrugsDetailResDto struct {
	Status  bool             `json:"status"`
	Message string           `json:"message"`
	Data    PharmacyDrugsRes `json:"data"`
}

type PharmacyDrugsRes struct {
	Id                   primitive.ObjectID  `json:"id" bson:"id"`
	Customer             CustomerInformation `json:"customer" bson:"customer"`
	FacilityOrProfession string              `json:"facilityOrProfession" bson:"facilityOrProfession"`
	PricePaid            float64             `json:"pricePaid" bson:"pricePaid"`
	ModeOfDelivery       string              `json:"modeOfDelivery" bson:"modeOfDelivery"`
	NameAndQuantity      string              `json:"nameAndQuantity" bson:"nameAndQuantity"`
	Prescription         []string            `json:"prescription" bson:"prescription"`
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

type AppointmentDetails struct {
	AppointmentFromDate time.Time `json:"appointmentFromDate" bson:"appointmentFromDate"`
	AppointmentToDate   time.Time `json:"appointmentToDate" bson:"appointmentToDate"`
}
