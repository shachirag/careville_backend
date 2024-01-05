package services

import "go.mongodb.org/mongo-driver/bson/primitive"

type LaboratoryRequestDto struct {
	LaboratoryReqDto LaboratoryReqDto `json:"data" form:"data"`
}

type LaboratoryReqDto struct {
	ProviderId           primitive.ObjectID `json:"providerId" form:"providerId"`
	Role                 string             `json:"role" form:"role"`
	FacilityOrProfession string             `json:"facilityOrProfession" form:"facilityOrProfession"`
	InformationName      string             `json:"informationName" form:"informationName"`
	Address              string             `json:"address" form:"address"`
	Longitude            string             `json:"longitude" form:"longitude"`
	Latitude             string             `json:"latitude" form:"latitude"`
	AdditionalText       string             `json:"additionalText" form:"additionalText"`
	Certificate          string             `json:"certificate" form:"certificate"`
	License              string             `json:"license" form:"license"`
	Investigations       []Investigations   `json:"investigations" form:"investigations"`
}

type Investigations struct {
	Type        string  `json:"type" bson:"type"`
	Name        string  `json:"name" bson:"name"`
	Information string  `json:"information" bson:"information"`
	Price       float64 `json:"price" bson:"price"`
}

type LaboratoryResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}
