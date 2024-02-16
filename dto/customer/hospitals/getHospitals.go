package hospitals

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetHospitalsRes struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Image     string             `json:"image" bson:"image"`
	Name      string             `json:"name" bson:"name"`
	Address   Address            `json:"address" bson:"address"`
	AvgRating float64            `json:"avgRating" bson:"avgRating"`
}

type Address struct {
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
	Add         string    `json:"add" bson:"add"`
	Type        string    `json:"type" bson:"type"`
}

type GetHospitalResDto struct {
	Status  bool              `json:"status"`
	Message string            `json:"message"`
	Data    []GetHospitalsRes `json:"data"`
}

type GetEmergencyHospitalResDto struct {
	Status  bool                       `json:"status"`
	Message string                     `json:"message"`
	Data    []GetEmergencyHospitalsRes `json:"data"`
}

type GetEmergencyHospitalsRes struct {
	Id                   primitive.ObjectID `json:"id" bson:"_id"`
	Image                string             `json:"image" bson:"image"`
	Name                 string             `json:"name" bson:"name"`
	Address              Address            `json:"address" bson:"address"`
	IsEmergencyAvailable bool               `json:"isEmergencyAvailable" bson:"isEmergencyAvailable"`
	Fees                 float64            `json:"fees" bson:"fees"`
}
