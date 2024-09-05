package services

import "go.mongodb.org/mongo-driver/bson/primitive"

type SendPharmacyDrugsInfoReqDto struct {
	AvailableDrugs     string             `json:"availableDrugs" bson:"availableDrugs"`
	NotAvailableDrugs  string             `json:"notAvailableDrugs" bson:"notAvailableDrugs"`
	TotalPriceToBePaid float64            `json:"totalPriceToBePaid" bson:"totalPriceToBePaid"`
	HomeDelivery       string             `json:"homeDelivery" bson:"homeDelivery"`
	DoctorApprovel     string             `json:"doctorApprovel" bson:"doctorApprovel"`
	AppointmentId      primitive.ObjectID `json:"appointmentId" bson:"appointmentId"`
}

type SendPharmacyDrugsInfoResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}

type GetPharmacyDrugsInfoResDto struct {
	Status  bool                 `json:"status" bson:"status"`
	Message string               `json:"message" bson:"message"`
	Data    PharmacyDrugsInfoRes `json:"data" bson:"data"`
}

type PharmacyDrugsInfoRes struct {
	AvailableDrugs     string             `json:"availableDrugs" bson:"availableDrugs"`
	NotAvailableDrugs  string             `json:"notAvailableDrugs" bson:"notAvailableDrugs"`
	TotalPriceToBePaid float64            `json:"totalPriceToBePaid" bson:"totalPriceToBePaid"`
	HomeDelivery       string             `json:"homeDelivery" bson:"homeDelivery"`
	DoctorApprovel     string             `json:"doctorApprovel" bson:"doctorApprovel"`
	AppointmentId      primitive.ObjectID `json:"appointmentId" bson:"appointmentId"`
}
