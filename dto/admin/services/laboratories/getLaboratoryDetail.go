package laboratory

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetLaboratoryDetailResDto struct {
	Status  bool                   `json:"status" bson:"status"`
	Message string                 `json:"message" bson:"message"`
	Data    GetLaboratoryDetailRes `json:"data" bson:"data"`
}

type GetLaboratoryDetailRes struct {
	Id                    primitive.ObjectID    `json:"id" bson:"id"`
	ProfileId             string                `json:"profileId" bson:"profileId"`
	FacilityOrProfession  string                `json:"facilityOrProfession" bson:"facilityOrProfession"`
	Role                  string                `json:"role" bson:"role"`
	User                  User                  `json:"user" bson:"user"`
	LaboratoryInformation LaboratoryInformation `json:"information" bson:"information"`
	Investigations        []Investigations      `json:"investigations" bson:"investigations"`
	Documents             Documents             `json:"documents" bson:"documents"`
	ServiceStatus         string                `json:"serviceStatus" bson:"serviceStatus"`
}

type User struct {
	FirstName   string      `json:"firstName" bson:"firstName"`
	LastName    string      `json:"lastName" bson:"lastName"`
	Email       string      `json:"email" bson:"email"`
	PhoneNumber PhoneNumber `json:"phoneNumber" bson:"phoneNumber"`
}

type Investigations struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	Type        string             `json:"type" bson:"type"`
	Information string             `json:"information" bson:"information"`
	Price       float64            `json:"price" bson:"price"`
}

type LaboratoryInformation struct {
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
