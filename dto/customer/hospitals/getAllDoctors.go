package hospitals

import "go.mongodb.org/mongo-driver/bson/primitive"

type DoctorResDto struct {
	Status  bool                   `json:"status" bson:"status"`
	Message string                 `json:"message" bson:"message"`
	Data    []SpecialityDoctorsRes `json:"data" bson:"data"`
}

type SpecialityDoctorsRes struct {
	Speciality string      `json:"speciality"`
	Doctors    []DoctorRes `json:"doctors"`
}

type DoctorRes struct {
	Id            primitive.ObjectID `json:"id"`
	Name          string             `json:"name"`
	Speciality    string             `json:"speciality"`
	Image         string             `json:"image"`
	NextAvailable NextAvailable      `json:"nextAvailable"`
}

type NextAvailable struct {
	StartTime string `json:"startTime" bson:"startTime"`
	// LastTime  string `json:"lastTime" bson:"lastTime"`
}
