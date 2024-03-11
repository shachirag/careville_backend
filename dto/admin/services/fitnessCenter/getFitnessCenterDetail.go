package fitnessCenter

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetFitnessCenterDetailResDto struct {
	Status  bool                      `json:"status" bson:"status"`
	Message string                    `json:"message" bson:"message"`
	Data    GetFitnessCenterDetailRes `json:"data" bson:"data"`
}

type GetFitnessCenterDetailRes struct {
	Id                       primitive.ObjectID       `json:"id" bson:"id"`
	ProfileId                string                   `json:"profileId" bson:"profileId"`
	FacilityOrProfession     string                   `json:"facilityOrProfession" bson:"facilityOrProfession"`
	Role                     string                   `json:"role" bson:"role"`
	User                     User                     `json:"user" bson:"user"`
	FitnessCenterInformation FitnessCenterInformation `json:"information" bson:"information"`
	Trainers                 []Trainers               `json:"trainers" bson:"trainers"`
	Subscripions             []Subscripions           `json:"subscripions" bson:"subscripions"`
	OtherServices            []OtherServices          `json:"otherServices" bson:"otherServices"`
	Documents                Documents                `json:"documents" bson:"documents"`
	ServiceStatus            string                   `json:"serviceStatus" bson:"serviceStatus"`
}

type User struct {
	FirstName   string      `json:"firstName" bson:"firstName"`
	LastName    string      `json:"lastName" bson:"lastName"`
	Email       string      `json:"email" bson:"email"`
	PhoneNumber PhoneNumber `json:"phoneNumber" bson:"phoneNumber"`
}

type Trainers struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	Category    string             `json:"category" bson:"category"`
	Information string             `json:"information" bson:"information"`
	Price       float64            `json:"price" bson:"price"`
}

type OtherServices struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	Information string             `json:"information" bson:"information"`
}

type Subscripions struct {
	Id      primitive.ObjectID `json:"id" bson:"id"`
	Type    string             `json:"type" bson:"type"`
	Details string             `json:"details" bson:"details"`
	Price   float64            `json:"price" bson:"price"`
}

type FitnessCenterInformation struct {
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
