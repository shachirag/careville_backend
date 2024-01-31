package services

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetMedicalLabScientistProfessionalDetailsResponseDto struct {
	Status  bool                                      `json:"status"`
	Message string                                    `json:"message"`
	Data    MedicalLabScientistProfessionalDetailsRes `json:"data"`
}

type MedicalLabScientistProfessionalDetailsRes struct {
	Qualification           string `json:"qualification"`
	Department              string `json:"department"`
	ProfessionalLicense     string `json:"professionalLicense"`
	ProfessionalCertificate string `json:"professionalCertificate"`
}

type MedicalLabScientistServicesResDto struct {
	Status  bool                            `json:"status" bson:"status"`
	Message string                          `json:"message" bson:"message"`
	Data    []MedicalLabScientistServiceRes `json:"data" bson:"data"`
}

type GetMedicalLabScientistServicesResDto struct {
	Status  bool                          `json:"status" bson:"status"`
	Message string                        `json:"message" bson:"message"`
	Data    MedicalLabScientistServiceRes `json:"data" bson:"data"`
}

type MedicalLabScientistServiceRes struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	ServiceFees float64            `json:"serviceFees" bson:"serviceFees"`
	Slots       []Slots            `json:"slots" bson:"slots"`
}

type UpdateMedicalLabScientistProfessionalInfoReqDto struct {
	Qualifications string `json:"qualifications" bson:"qualifications"`
	Department     string `json:"department" bson:"department"`
}

type UpdateMedicalLabScientistProfessionalInfoResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}

type MoreMedicalLabScientistServiceReqDto struct {
	Name        string  `json:"name" bson:"name"`
	ServiceFees float64 `json:"serviceFees" bson:"serviceFees"`
	Slots       []Slots `json:"slots" bson:"slots"`
}

type MoreMedicalLabScientistServiceResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}

type UpdateMedicalLabScientistServiceReqDto struct {
	Name        string  `json:"name" bson:"name"`
	ServiceFees float64 `json:"serviceFees" bson:"serviceFees"`
	Slots       []Slots `json:"slots" bson:"slots"`
}

type UpdateMedicalLabScientistServiceResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}
