package services

import "go.mongodb.org/mongo-driver/bson/primitive"

type DoctorProfessionRequestDto struct {
	DoctorProfessionReqDto DoctorProfessionReqDto `json:"data" form:"data"`
}

type DoctorProfessionReqDto struct {
	ProviderId                 primitive.ObjectID         `json:"providerId" form:"providerId"`
	Role                       string                     `json:"role" form:"role"`
	FacilityOrProfession       string                     `json:"facilityOrProfession" form:"facilityOrProfession"`
	InformationName            string                     `json:"informationName" form:"informationName"`
	Address                    string                     `json:"address" form:"address"`
	Longitude                  string                     `json:"longitude" form:"longitude"`
	Latitude                   string                     `json:"latitude" form:"latitude"`
	AdditionalText             string                     `json:"additionalText" form:"additionalText"`
	Certificate                string                     `json:"certificate" form:"certificate"`
	License                    string                     `json:"license" form:"license"`
	Speciality                 string                     `json:"speciality" bson:"speciality"`
	Qualifications             string                     `json:"qualifications" bson:"qualifications"`
	PersonalIdentificationDocs PersonalIdentificationDocs `json:"personalIdentificationDocs" form:"personalIdentificationDocs"`
	ProfessionalDetailsDocs    ProfessionalDetailsDocs    `json:"professionalDetailsDocs" form:"professionalDetailsDocs"`
	Schedule                   []DoctorSchedule           `json:"schedule" form:"schedule"`
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
