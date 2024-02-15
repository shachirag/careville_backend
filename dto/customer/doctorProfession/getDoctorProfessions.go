package doctorProfession

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetDoctorProfessionRes struct {
	Id            primitive.ObjectID `json:"id" bson:"_id"`
	Image         string             `json:"image" bson:"image"`
	Name          string             `json:"name" bson:"name"`
	Speciality    string             `json:"speciality" bson:"speciality"`
	NextAvailable NextAvailable      `json:"nextAvailable" bson:"nextAvailable"`
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
