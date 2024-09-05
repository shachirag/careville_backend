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
	Id                         primitive.ObjectID         `json:"id" bson:"id"`
	FacilityOrProfession       string                     `json:"facilityOrProfession" bson:"facilityOrProfession"`
	PharmacyInformation        PharmacyInformation        `json:"pharmacyInformation" bson:"pharmacyInformation"`
	ProvderProvidedInformation ProvderProvidedInformation `json:"provderProvidedInformation" bson:"provderProvidedInformation"`
	// Customer             CustomerInformation `json:"customer" bson:"customer"`
	// PricePaid            float64             `json:"pricePaid" bson:"pricePaid"`
	// ModeOfDelivery       string              `json:"modeOfDelivery" bson:"modeOfDelivery"`
	// NameAndQuantity      string              `json:"nameAndQuantity" bson:"nameAndQuantity"`
	// Prescription         []string            `json:"prescription" bson:"prescription"`
}

type ProvderProvidedInformation struct {
	AvailableDrugs     string  `json:"availableDrugs" bson:"availableDrugs"`
	NotAvailableDrugs  string  `json:"notAvailableDrugs" bson:"notAvailableDrugs"`
	DoctorApprovel     string  `json:"doctorApprovel" bson:"doctorApprovel"`
	HomeDelivery       string  `json:"homeDelivery" bson:"homeDelivery"`
	TotalPriceToBePaid float64 `json:"totalPriceToBePaid" bson:"totalPriceToBePaid"`
}

type PharmacyInformation struct {
	Id        primitive.ObjectID `json:"id" bson:"id"`
	Address   Address            `json:"address" bson:"address"`
	Image     string             `json:"image" bson:"image"`
	Name      string             `json:"name" bson:"name"`
	AvgRating float64            `json:"avgRating" bson:"avgRating"`
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
