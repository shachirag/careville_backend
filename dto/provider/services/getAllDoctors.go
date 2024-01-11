package services

import "go.mongodb.org/mongo-driver/bson/primitive"

type DoctorResDto struct {
	Status  bool        `json:"status" bson:"status"`
	Message string      `json:"message" bson:"message"`
	Data    []DoctorRes `json:"data" bson:"data"`
}

type DoctorRes struct {
	Id         primitive.ObjectID  `json:"id"`
	Name       string              `json:"name"`
	Speciality string              `json:"speciality"`
	Schedule   []DoctorScheduleRes `json:"schedule"`
}

type DoctorScheduleRes struct {
	StartTime string   `json:"startTime" bson:"startTime"`
	EndTime   string   `json:"endTime" bson:"endTime"`
	Days      []string `json:"days" bson:"days"`
}

type GetDoctorResDto struct {
	Status  bool         `json:"status" bson:"status"`
	Message string       `json:"message" bson:"message"`
	Data    DoctorRes `json:"data" bson:"data"`
}
