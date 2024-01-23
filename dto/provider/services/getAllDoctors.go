package services

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
	Id         primitive.ObjectID  `json:"id"`
	Name       string              `json:"name"`
	Speciality string              `json:"speciality"`
	Image      string              `json:"image"`
	Schedule   []DoctorScheduleRes `json:"schedule"`
}

type DoctorScheduleRes struct {
	StartTime string   `json:"startTime" bson:"startTime"`
	EndTime   string   `json:"endTime" bson:"endTime"`
	Days      []string `json:"days" bson:"days"`
}

type GetDoctorResDto struct {
	Status  bool      `json:"status" bson:"status"`
	Message string    `json:"message" bson:"message"`
	Data    DoctorRes `json:"data" bson:"data"`
}

type OtherServicesResDto struct {
	Status  bool             `json:"status" bson:"status"`
	Message string           `json:"message" bson:"message"`
	Data    OtherServicesRes `json:"data" bson:"data"`
}

type OtherServicesRes struct {
	OtherServices []string `json:"otherServices" bson:"otherServices"`
}

type NotificationResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}
