package doctorProfession

type AddEmergencyDoctorReqDto struct {
	DoctorId  string  `json:"doctorId" form:"doctorId"`
	PricePaid float64 `json:"pricePaid" form:"pricePaid"`
}

type AddEmergencyHospitalReqDto struct {
	HospitalId string  `json:"hospitalId" form:"hospitalId"`
	Address    string  `json:"address" form:"address"`
	Longitude  string  `json:"longitude" form:"longitude"`
	Latitude   string  `json:"latitude" form:"latitude"`
}

type AddEmergencyDoctorResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
