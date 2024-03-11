package physiotherapist

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetPhysiotherapistDetailResDto struct {
	Status  bool                        `json:"status" bson:"status"`
	Message string                      `json:"message" bson:"message"`
	Data    GetPhysiotherapistDetailRes `json:"data" bson:"data"`
}

type GetPhysiotherapistDetailRes struct {
	Id                         primitive.ObjectID         `json:"id" bson:"id"`
	ProfileId                  string                     `json:"profileId" bson:"profileId"`
	FacilityOrProfession       string                     `json:"facilityOrProfession" bson:"facilityOrProfession"`
	Role                       string                     `json:"role" bson:"role"`
	User                       User                       `json:"user" bson:"user"`
	PhysiotherapistInformation PhysiotherapistInformation `json:"information" bson:"information"`
	ProfessionalDetails        ProfessionalDetails        `json:"professionalDetails" bson:"professionalDetails"`
	ProfessionalDocuments      ProfessionalDocuments      `json:"professionalDocuments" bson:"professionalDocuments"`
	ServiceAndSchedule         []ServiceAndSchedule       `json:"serviceAndSchedule" bson:"serviceAndSchedule"`
	PersonalDocuments          PersonalDocuments          `json:"personalDocuments" bson:"personalDocuments"`
	ServiceStatus              string                     `json:"serviceStatus" bson:"serviceStatus"`
}

type User struct {
	FirstName   string      `json:"firstName" bson:"firstName"`
	LastName    string      `json:"lastName" bson:"lastName"`
	Email       string      `json:"email" bson:"email"`
	PhoneNumber PhoneNumber `json:"phoneNumber" bson:"phoneNumber"`
}

type ProfessionalDetails struct {
	Qualification string `json:"qualification" bson:"qualification"`
}

type PhysiotherapistInformation struct {
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

type ProfessionalDocuments struct {
	Certificate string `json:"certificate" bson:"certificate"`
	License     string `json:"license" bson:"license"`
}

type PersonalDocuments struct {
	Nimc    string `json:"nimc" bson:"nimc"`
	License string `json:"license" bson:"license"`
}

type ServiceAndSchedule struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	ServiceFees float64            `json:"serviceFees" bson:"serviceFees"`
	Slots       []Slots            `json:"slots" bson:"slots"`
}

type Slots struct {
	StartTime string   `json:"startTime" bson:"startTime"`
	EndTime   string   `json:"endTime" bson:"endTime"`
	Days      []string `json:"days" bson:"days"`
}
