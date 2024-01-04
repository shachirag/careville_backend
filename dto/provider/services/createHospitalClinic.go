package hospitalclinic

import "go.mongodb.org/mongo-driver/bson/primitive"

type HospitalClinicReqDto struct {
	ProviderId     primitive.ObjectID `json:"providerId" form:"providerId"`
	Name           string             `json:"name" form:"name"`
	Address        string             `json:"address" form:"address"`
	Longitude      string             `json:"longitude" form:"longitude"`
	Latitude       string             `json:"latitude" form:"latitude"`
	AdditionalText string             `json:"additionalText" form:"additionalText"`
	OtherServices  []string           `json:"otherServices" form:"otherServices"`
	Insurances     []string           `json:"insurances" form:"insurances"`
	HopitalImage   string             `json:"hospitalImage" form:"insurances"`
	Certificate    string             `json:"certificate" form:"certificate"`
	License        string             `json:"license" form:"license"`
	Doctor         []Doctor           `json:"doctor" form:"doctor"`
}

type Doctor struct {
	Name       string     `json:"name" form:"name"`
	Speciality string     `json:"speciality" form:"speciality"`
	Schedule   []Schedule `json:"schedule" form:"schedule"`
}

type Schedule struct {
	StartTime string   `json:"startTime" form:"startTime"`
	EndTime   string   `json:"endTime" form:"endTime"`
	Days      []string `json:"days" form:"days"`
}
