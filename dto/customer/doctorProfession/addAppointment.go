package doctorProfession

type DoctorProfessionAppointmentReqDto struct {
	FamillyMemberId string  `json:"familyMemberId" form:"familyMemberId"`
	AppointmentDate string  `json:"appointmentDate" form:"appointmentDate"`
	AvailableTime   string  `json:"availableTime" form:"availableTime"`
	RemindMeBefore  string  `json:"remindMeBefore" form:"remindMeBefore"`
	FamilyType      string  `json:"familyType" form:"familyType"`
	PricePaid       float64 `json:"pricePaid" form:"pricePaid"`
	Address         string  `json:"address" form:"address"`
	Longitude       string  `json:"longitude" form:"longitude"`
	Latitude        string  `json:"latitude" form:"latitude"`
}

type DoctorProfessionAppointmentResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
