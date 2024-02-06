package nurse

type NurseAppointmentReqDto struct {
	FamillyMemberId string  `json:"familyMemberId" form:"familyMemberId"`
	AppointmentDate string  `json:"appointmentDate" form:"appointmentDate"`
	AvailableTime   string  `json:"availableTime" form:"availableTime"`
	RemindMeBefore  string  `json:"remindMeBefore" form:"remindMeBefore"`
	FamilyType      string  `json:"familyType" form:"familyType"`
	PricePaid       float64 `json:"pricePaid" form:"pricePaid"`
}

type NurseAppointmentResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
