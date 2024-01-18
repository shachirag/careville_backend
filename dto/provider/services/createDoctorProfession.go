package services

import "go.mongodb.org/mongo-driver/bson/primitive"

type DoctorProfessionRequestDto struct {
	DoctorProfessionReqDto DoctorProfessionReqDto `json:"data" form:"data"`
}

type DoctorProfessionReqDto struct {
	InformationFirstName string           `json:"informationFirstName" form:"informationFirstName"`
	InformationLastName  string           `json:"informationLastName" form:"informationLastName"`
	Address              string           `json:"address" form:"address"`
	Longitude            string           `json:"longitude" form:"longitude"`
	Latitude             string           `json:"latitude" form:"latitude"`
	AdditionalText       string           `json:"additionalText" form:"additionalText"`
	Certificate          string           `json:"certificate" form:"certificate"`
	Speciality           string           `json:"speciality" bson:"speciality"`
	Qualifications       string           `json:"qualifications" bson:"qualifications"`
	Schedule             []DoctorSchedule `json:"schedule" form:"schedule"`
}

type PersonalIdentificationDocs struct {
	Nimc    string `json:"nimc" bson:"nimc"`
	License string `json:"license" bson:"license"`
}

type ProfessionalDetailsDocs struct {
	Certificate string `json:"certificate" bson:"certificate"`
	License     string `json:"license" bson:"license"`
}

type DoctorSchedule struct {
	ConsultationFees string `json:"consultationFees" bson:"consultationFees"`
	Slots            Slots  `json:"slots" bson:"slots"`
}

type Slots struct {
	StartTime string   `json:"startTime" bson:"startTime"`
	EndTime   string   `json:"endTime" bson:"endTime"`
	Days      []string `json:"days" bson:"days"`
}

type DoctorProfessionResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}

type StatusResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}

type StatusRes struct {
	Status  bool          `json:"status"`
	Message string        `json:"message"`
	Data    StatusRespDto `json:"data"`
}

type StatusRespDto struct {
	Id     primitive.ObjectID `json:"id" bson:"_id"`
	Status string             `json:"status" bson:"status"`
}
