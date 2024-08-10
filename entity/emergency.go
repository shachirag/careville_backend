package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EmergencyEntity struct {
	ID                   primitive.ObjectID       `json:"id" bson:"_id"`
	Doctor               *DoctorEmergencyEntity   `json:"doctor,omitempty" bson:"doctor,omitempty"`
	Hospital             *HospitalEmergencyEntity `json:"hospital,omitempty" bson:"hospital,omitempty"`
	Price                PriceEmergencyEntity     `json:"price" bson:"price"`
	Type                 string                   `json:"type" bson:"type"`
	Customer             CustomerEmergencyEntity  `json:"customer" bson:"customer"`
	CreatedAt            time.Time                `json:"createdAt" bson:"createdAt"`
	UpdatedAt            time.Time                `json:"updatedAt" bson:"updatedAt"`
	ServiceID            primitive.ObjectID       `json:"serviceId" bson:"serviceId"`
	FacilityOrProfession string                   `json:"facilityOrProfession" bson:"facilityOrProfession"`
	Role                 string                   `json:"role" bson:"role"`
}

type DoctorEmergencyEntity struct {
	ID         primitive.ObjectID `json:"id" bson:"id"`
	Name       string             `json:"name" bson:"name"`
	Speciality string             `json:"speciality" bson:"speciality"`
	Image      string             `json:"image" bson:"image"`
}

type HospitalEmergencyEntity struct {
	Information  HospitalInformation `json:"information" bson:"information"`
	AddedAddress Address             `json:"addedAddress" bson:"addedAddress"`
}

type HospitalInformation struct {
	ID      primitive.ObjectID `json:"id" bson:"id"`
	Name    string             `json:"name" bson:"name"`
	Address Address            `json:"address" bson:"address"`
	Image   string             `json:"image" bson:"image"`
}

type CustomerEmergencyEntity struct {
	ID          primitive.ObjectID `json:"id" bson:"id"`
	FirstName   string             `json:"firstName" bson:"firstName"`
	LastName    string             `json:"lastName" bson:"lastName"`
	Image       string             `json:"image" bson:"image"`
	Email       string             `json:"email" bson:"email"`
	PhoneNumber PhoneNumber        `json:"phoneNumber" bson:"phoneNumber"`
	Address     Address            `json:"address" bson:"address"`
}

type PriceEmergencyEntity struct {
	PricePaid float64 `json:"pricePaid" bson:"pricePaid"`
}
