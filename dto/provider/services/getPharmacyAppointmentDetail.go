package services

import (
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

type Address struct {
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
	Type        string    `json:"type" bson:"type"`
	Add         string    `json:"add" bson:"add"`
}
