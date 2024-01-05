package services

import "go.mongodb.org/mongo-driver/bson/primitive"

type PharmacyRequestDto struct {
	PharmacyReqDto PharmacyReqDto `json:"data" form:"data"`
}

type PharmacyReqDto struct {
	ProviderId           primitive.ObjectID   `json:"providerId" form:"providerId"`
	Role                 string               `json:"role" form:"role"`
	FacilityOrProfession string               `json:"facilityOrProfession" form:"facilityOrProfession"`
	InformationName      string               `json:"informationName" form:"informationName"`
	Address              string               `json:"address" form:"address"`
	Longitude            string               `json:"longitude" form:"longitude"`
	Latitude             string               `json:"latitude" form:"latitude"`
	AdditionalText       string               `json:"additionalText" form:"additionalText"`
	Certificate          string               `json:"certificate" form:"certificate"`
	License              string               `json:"license" form:"license"`
	AdditionalServices   []AdditionalServices `json:"additionalServices" form:"additionalServices"`
}

type PharmacyResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}
