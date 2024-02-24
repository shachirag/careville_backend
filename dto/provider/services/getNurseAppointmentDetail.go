package services

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetNurseAppointmentDetailResDto struct {
	Status  bool                   `json:"status"`
	Message string                 `json:"message"`
	Data    NurseAppointmentRes `json:"data"`
}

type NurseAppointmentRes struct {
	Id                   primitive.ObjectID  `json:"id" bson:"id"`
	Customer             CustomerInformation `json:"customer" bson:"customer"`
	AppointmentDetails   AppointmentDetails  `json:"appointmentDetails" bson:"appointmentDetails"`
	FacilityOrProfession string              `json:"facilityOrProfession" bson:"facilityOrProfession"`
	PricePaid            float64             `json:"pricePaid" bson:"pricePaid"`
	FamilyMember         FamilyMember        `json:"familyMember"`
}
