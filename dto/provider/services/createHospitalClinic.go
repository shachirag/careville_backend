package services

type HospitalClinicRequestDto struct {
	HospitalClinicReqDto HospitalClinicReqDto `json:"data" form:"data"`
}

type HospitalClinicReqDto struct {
	InformationName string   `json:"informationName" form:"informationName"`
	Address         string   `json:"address" form:"address"`
	Longitude       string   `json:"longitude" form:"longitude"`
	Latitude        string   `json:"latitude" form:"latitude"`
	AdditionalText  string   `json:"additionalText" form:"additionalText"`
	OtherServices   []string `json:"otherServices" form:"otherServices"`
	Insurances      []string `json:"insurances" form:"insurances"`
	Doctor          []Doctor `json:"doctor" form:"doctor"`
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

type HospitalClinicResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}
