package services

type MedicalLabScientistRequestDto struct {
	MedicalLabScientistReqDto MedicalLabScientistReqDto `json:"data" form:"data"`
}

type MedicalLabScientistReqDto struct {
	InformationFirstName string                    `json:"informationFirstName" form:"informationFirstName"`
	InformationLastName  string                    `json:"informationLastName" form:"informationLastName"`
	Address              string                    `json:"address" form:"address"`
	Longitude            string                    `json:"longitude" form:"longitude"`
	Latitude             string                    `json:"latitude" form:"latitude"`
	AdditionalText       string                    `json:"additionalText" form:"additionalText"`
	Department           string                    `json:"department" form:"department"`
	Document             string                    `json:"document" form:"document"`
	Schedule             []PhysiotherapistSchedule `json:"schedule" form:"schedule"`
}

type MedicalLabScientistSchedule struct {
	Name        string `json:"name" bson:"name"`
	ServiceFees string `json:"serviceFees" bson:"serviceFees"`
	Slots       Slots  `json:"slots" bson:"slots"`
}

type MedicalLabScientistResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}
