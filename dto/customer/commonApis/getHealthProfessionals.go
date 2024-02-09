package common

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetHealthProfessionalResDto struct {
	Status  bool                     `json:"status"`
	Message string                   `json:"message"`
	Data    HealthProfessionalResDto `json:"data"`
}

type HealthProfessionalResDto struct {
	Doctors              []GetDoctorHealthProfessionalRes `json:"doctors"`
	Physiotherapists     []GetHealthProfessionalRes       `json:"physiotherapists"`
	MedicalLabScientists []GetHealthProfessionalRes       `json:"medicalLabScientists"`
	Nurse                []GetHealthProfessionalRes       `json:"nurses"`
}

type GetDoctorHealthProfessionalRes struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	Image      string             `json:"image" bson:"image"`
	Name       string             `json:"name" bson:"name"`
	AvgRating  float64            `json:"avgRating" bson:"avgRating"`
	Speciality string             `json:"speciality" bson:"speciality"`
}

type GetHealthProfessionalRes struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Image     string             `json:"image" bson:"image"`
	Name      string             `json:"name" bson:"name"`
	AvgRating float64            `json:"avgRating" bson:"avgRating"`
}
