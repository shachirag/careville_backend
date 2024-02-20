package doctorProfession

type AddEmergencyDoctorReqDto struct {
	DoctorId  string  `json:"doctorId" form:"doctorId"`
	PricePaid float64 `json:"pricePaid" form:"pricePaid"`
}

type AddEmergencyDoctorResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
