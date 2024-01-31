package services

import "go.mongodb.org/mongo-driver/bson/primitive"

type DoctorProfessionRequestDto struct {
	DoctorProfessionReqDto DoctorProfessionReqDto `json:"data" form:"data"`
}

type DoctorProfessionReqDto struct {
	InformationName  string  `json:"informationName" form:"informationName"`
	Address          string  `json:"address" form:"address"`
	Longitude        string  `json:"longitude" form:"longitude"`
	Latitude         string  `json:"latitude" form:"latitude"`
	AdditionalText   string  `json:"additionalText" form:"additionalText"`
	Certificate      string  `json:"certificate" form:"certificate"`
	Speciality       string  `json:"speciality" bson:"speciality"`
	Qualifications   string  `json:"qualifications" bson:"qualifications"`
	Slots            []Slots `json:"slots" bson:"slots"`
	ConsultationFees float64 `json:"consultationFees" bson:"consultationFees"`
	// Schedule         []DoctorSchedule `json:"schedule" form:"schedule"`
}

type PersonalIdentificationDocs struct {
	Nimc    string `json:"nimc" bson:"nimc"`
	License string `json:"license" bson:"license"`
}

type ProfessionalDetailsDocs struct {
	Certificate string `json:"certificate" bson:"certificate"`
	License     string `json:"license" bson:"license"`
}

type Slots struct {
	StartTime string   `json:"startTime" bson:"startTime"`
	EndTime   string   `json:"endTime" bson:"endTime"`
	Days      []string `json:"days" bson:"days"`
}

type DoctorProfessionResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
	Role    Role   `json:"data" bson:"data"`
}

type Role struct {
	Role                 string `json:"role" bson:"role"`
	FacilityOrProfession string `json:"facilityOrProfession" bson:"facilityOrProfession"`
	ServiceStatus        string `json:"serviceStatus" bson:"serviceStatus"`
	Name                 string `json:"name" bson:"name"`
	Image                string `json:"image" bson:"image"`
	IsEmergencyAvailable bool   `json:"isEmergencyAvailable" bson:"isEmergencyAvailable"`
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
	Id            primitive.ObjectID `json:"id" bson:"_id"`
	ServiceStatus string             `json:"serviceStatus" bson:"serviceStatus"`
}
