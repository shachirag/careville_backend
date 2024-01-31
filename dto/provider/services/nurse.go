package services

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetNurseProfessionalDetailsResDto struct {
	Status  bool                      `json:"status"`
	Message string                    `json:"message"`
	Data    PhysiotherapistDetailsRes `json:"data"`
}

type NurseDetailsRes struct {
	Qualification           string `json:"qualification"`
	ProfessionalLicense     string `json:"professionalLicense"`
	ProfessionalCertificate string `json:"professionalCertificate"`
}

type UpdateNurseProfessionalInfoReqDto struct {
	Qualifications string `json:"qualifications" bson:"qualifications"`
}

type UpdateNurseProfessionalInfoResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}

type NurseServicesResDto struct {
	Status  bool              `json:"status" bson:"status"`
	Message string            `json:"message" bson:"message"`
	Data    []NurseServiceRes `json:"data" bson:"data"`
}

type GetNurseServicesResDto struct {
	Status  bool            `json:"status" bson:"status"`
	Message string          `json:"message" bson:"message"`
	Data    NurseServiceRes `json:"data" bson:"data"`
}

type NurseServiceRes struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	ServiceFees float64            `json:"serviceFees" bson:"serviceFees"`
	Slots       []Slots            `json:"slots" bson:"slots"`
}

type MoreNurseServiceReqDto struct {
	Name        string  `json:"name" bson:"name"`
	ServiceFees float64 `json:"serviceFees" bson:"serviceFees"`
	Slots       []Slots `json:"slots" bson:"slots"`
}

type MoreNurseServiceResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}

type DeleteNurseProfessionalInfoResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}

type UpdateNurseServiceReqDto struct {
	Name        string  `json:"name" form:"name"`
	ServiceFees float64 `json:"serviceFees" form:"serviceFees"`
	Slots       []Slots `json:"slots" form:"slots"`
}
