package services

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetLaboratoryAppointmentDetailResDto struct {
	Status  bool                   `json:"status"`
	Message string                 `json:"message"`
	Data    LaboratoryAppointmentRes `json:"data"`
}

type LaboratoryAppointmentRes struct {
	Id                   primitive.ObjectID  `json:"id" bson:"id"`
	Customer             CustomerInformation `json:"customer" bson:"customer"`
	AppointmentDetails   AppointmentData     `json:"appointmentDetails" bson:"appointmentDetails"`
	Investigation        Investigation       `json:"investigation" bson:"investigation"`
	FacilityOrProfession string              `json:"facilityOrProfession" bson:"facilityOrProfession"`
	PricePaid            float64             `json:"pricePaid" bson:"pricePaid"`
	FamilyMember         FamilyMember        `json:"familyMember"`
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
