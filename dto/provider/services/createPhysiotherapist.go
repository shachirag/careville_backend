package services

import "go.mongodb.org/mongo-driver/bson/primitive"

type PhysiotherapistRequestDto struct {
	PhysiotherapistReqDto PhysiotherapistReqDto `json:"data" form:"data"`
}

type PhysiotherapistReqDto struct {
	ProviderId           primitive.ObjectID        `json:"providerId" form:"providerId"`
	Role                 string                    `json:"role" form:"role"`
	FacilityOrProfession string                    `json:"facilityOrProfession" form:"facilityOrProfession"`
	InformationName      string                    `json:"informationName" form:"informationName"`
	Address              string                    `json:"address" form:"address"`
	Longitude            string                    `json:"longitude" form:"longitude"`
	Latitude             string                    `json:"latitude" form:"latitude"`
	AdditionalText       string                    `json:"additionalText" form:"additionalText"`
	Qualifications       string                    `json:"qualifications" bson:"qualifications"`
	Schedule             []PhysiotherapistSchedule `json:"schedule" form:"schedule"`
}

type PhysiotherapistSchedule struct {
	Name        string `json:"name" bson:"name"`
	ServiceFees string `json:"serviceFees" bson:"serviceFees"`
	Slots       Slots  `json:"slots" bson:"slots"`
}

type PhysiotherapistResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}
