package services

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetHospitalDetailResDto struct {
	Status  bool                 `json:"status" bson:"status"`
	Message string               `json:"message" bson:"message"`
	Data    GetHospitalDetailRes `json:"data" bson:"data"`
}

type GetHospitalDetailRes struct {
	Id                   primitive.ObjectID  `json:"id" bson:"id"`
	ProfileId            string              `json:"profileId" bson:"profileId"`
	FacilityOrProfession string              `json:"facilityOrProfession" bson:"facilityOrProfession"`
	Role                 string              `json:"role" bson:"role"`
	User                 User                `json:"user" bson:"user"`
	HospitalInformation  HospitalInformation `json:"information" bson:"information"`
	Doctor               []Doctor            `json:"doctor" bson:"doctor"`
	OtherServices        []string            `json:"otherServices" bson:"otherServices"`
	Insurances           []string            `json:"insurances" bson:"insurances"`
	Documents            Documents           `json:"documents" bson:"documents"`
	ServiceStatus        string              `json:"serviceStatus" bson:"serviceStatus"`
}

type User struct {
	FirstName   string      `json:"firstName" bson:"firstName"`
	LastName    string      `json:"lastName" bson:"lastName"`
	Email       string      `json:"email" bson:"email"`
	PhoneNumber PhoneNumber `json:"phoneNumber" bson:"phoneNumber"`
}

type Doctor struct {
	Id         primitive.ObjectID `json:"id" bson:"id"`
	Name       string             `json:"name" bson:"name"`
	Speciality string             `json:"speciality" bson:"speciality"`
	Schedule   []Schedule         `json:"schedule" bson:"schedule"`
}

type Schedule struct {
	StartTime string   `json:"startTime" bson:"startTime"`
	EndTime   string   `json:"endTime" bson:"endTime"`
	Days      []string `json:"days" bson:"days"`
}

type HospitalInformation struct {
	Name           string  `json:"name" bson:"name"`
	AdditionalText string  `json:"additionalText" bson:"additionalText"`
	Image          string  `json:"image" bson:"image"`
	Address        Address `json:"address" bson:"address"`
}

type Address struct {
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
	Add         string    `json:"add" bson:"add"`
	Type        string    `json:"type" bson:"type"`
}

type Documents struct {
	Certificate string `json:"certificate" bson:"certificate"`
	License     string `json:"license" bson:"license"`
}
