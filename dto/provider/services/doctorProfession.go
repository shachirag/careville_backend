package services

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetDoctorProfessionProfessionalDetailsResDto struct {
	Status  bool                                 `json:"status"`
	Message string                               `json:"message"`
	Data    DoctorProfessionProfessionDetailsRes `json:"data"`
}

type DoctorProfessionProfessionDetailsRes struct {
	Qualification           string  `json:"qualification"`
	ProfessionalLicense     string  `json:"professionalLicense"`
	ProfessionalCertificate string  `json:"professionalCertificate"`
	ConsultingFees          float64 `json:"consultingFees"`
	Speciality              string  `json:"speciality"`
}

type UpdateDoctorProfessionProfessionalInfoReqDto struct {
	Qualifications string  `json:"qualifications" bson:"qualifications"`
	ConsultingFees float64 `json:"consultingFees" bson:"consultingFees"`
	Speciality     string  `json:"speciality" bson:"speciality"`
}

type UpdateDoctorProfessionProfessionalInfoResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}

type DoctorProfessionSlotsResDto struct {
	Status  bool                               `json:"status" bson:"status"`
	Message string                             `json:"message" bson:"message"`
	Data    []DoctorSlotsResponseDto `json:"data" bson:"data"`
}

type GetDoctorProfessionSlotsResDto struct {
	Status  bool                   `json:"status" bson:"status"`
	Message string                 `json:"message" bson:"message"`
	Data    DoctorProfessionSlotsResponseDto `json:"data" bson:"data"`
}

type DoctorSlotsResponseDto struct {
	Id        primitive.ObjectID `json:"id" bson:"id"`
	StartTime string             `json:"startTime" bson:"startTime"`
	EndTime   string             `json:"endTime" bson:"endTime"`
	Days      []string           `json:"days" bson:"days"`
}

type DoctorProfessionSlotsResponseDto struct {
	Slots []DoctorSlots `json:"slots" bson:"slots"`
}

type DoctorSlots struct {
	Id        primitive.ObjectID `json:"id" bson:"id"`
	StartTime string             `json:"startTime" bson:"startTime"`
	EndTime   string             `json:"endTime" bson:"endTime"`
	Days      []string           `json:"days" bson:"days"`
}

type UpdateDoctorProfessionSlotReqDto struct {
	Slots []Slots `json:"slots" form:"slots"`
}

type EditDoctorProfessionSlotReqDto struct {
	StartTime string   `json:"startTime" bson:"startTime"`
	EndTime   string   `json:"endTime" bson:"endTime"`
	Days      []string `json:"days" bson:"days"`
}

type UpdateDoctorProfessionSlotResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}
