package services

type NurseRequestDto struct {
	NurseReqDto NurseReqDto `json:"data" form:"data"`
}

type NurseReqDto struct {
	InformationName string          `json:"informationName" form:"informationName"`
	Address         string          `json:"address" form:"address"`
	Longitude       string          `json:"longitude" form:"longitude"`
	Latitude        string          `json:"latitude" form:"latitude"`
	AdditionalText  string          `json:"additionalText" form:"additionalText"`
	Qualifications  string          `json:"qualifications" bson:"qualifications"`
	Schedule        []NurseSchedule `json:"schedule" form:"schedule"`
}

type NurseSchedule struct {
	Name        string `json:"name" bson:"name"`
	ServiceFees string `json:"serviceFees" bson:"serviceFees"`
	Slots       Slots  `json:"slots" bson:"slots"`
}

type NurseResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}
