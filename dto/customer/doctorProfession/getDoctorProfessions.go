package doctorProfession

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetDoctorProfessionRes struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Image       string             `json:"image" bson:"image"`
	Name        string             `json:"name" bson:"name"`
	Speciality  string             `json:"speciality" bson:"speciality"`
	ServiceType string             `json:"serviceType" bson:"serviceType"`
	NextAvailable NextAvailable `json:"nextAvailable" bson:"nextAvailable"`
}

type NextAvailable struct {
	StartTime string `json:"startTime" bson:"startTime"`
	LastTime  string `json:"lastTime" bson:"lastTime"`
}

type GetDoctorProfessionResponseDto struct {
	Status  bool                     `json:"status"`
	Message string                   `json:"message"`
	Data    []GetDoctorProfessionRes `json:"data"`
}

type GetEmergencyDoctorProfessionRes struct {
	Id                   primitive.ObjectID `json:"id" bson:"_id"`
	Image                string             `json:"image" bson:"image"`
	Name                 string             `json:"name" bson:"name"`
	Speciality           string             `json:"speciality" bson:"speciality"`
	IsEmergencyAvailable bool               `json:"isEmergencyAvailable" bson:"isEmergencyAvailable"`
	ConsultationFees     float64            `json:"consultationFees" bson:"consultationFees"`
}

type EmergencyDoctorResDto struct {
	Status  bool                              `json:"status"`
	Message string                            `json:"message"`
	Data    []GetEmergencyDoctorProfessionRes `json:"data"`
}
