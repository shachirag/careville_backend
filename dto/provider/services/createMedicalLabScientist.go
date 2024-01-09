package services

import "go.mongodb.org/mongo-driver/bson/primitive"

type MedicalLabScientistRequestDto struct {
	MedicalLabScientistReqDto MedicalLabScientistReqDto `json:"data" form:"data"`
}

type MedicalLabScientistReqDto struct {
	ProviderId           primitive.ObjectID        `json:"providerId" form:"providerId"`
	Role                 string                    `json:"role" form:"role"`
	FacilityOrProfession string                    `json:"facilityOrProfession" form:"facilityOrProfession"`
	InformationName      string                    `json:"informationName" form:"informationName"`
	Address              string                    `json:"address" form:"address"`
	Longitude            string                    `json:"longitude" form:"longitude"`
	Latitude             string                    `json:"latitude" form:"latitude"`
	AdditionalText       string                    `json:"additionalText" form:"additionalText"`
	Department           string                    `json:"department" form:"department"`
	Document             string                    `json:"document" form:"document"`
	Schedule             []PhysiotherapistSchedule `json:"schedule" form:"schedule"`
}

type MedicalLabScientistSchedule struct {
	Name        string `json:"name" bson:"name"`
	ServiceFees string `json:"serviceFees" bson:"serviceFees"`
	Slots       Slots  `json:"slots" bson:"slots"`
}

type MedicalLabScientistResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}
