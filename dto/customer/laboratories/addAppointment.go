package laboratory

type LaboratoryAppointmentReqDto struct {
	FamillyMemberId string  `json:"familyMemberId" form:"familyMemberId"`
	InvestigationId string  `json:"investigationId" form:"investigationId"`
	AppointmentDate string  `json:"appointmentDate" form:"appointmentDate"`
	FamilyType      string  `json:"familyType" form:"familyType"`
	PricePaid       float64 `json:"pricePaid" form:"pricePaid"`
}

type LaboratoryAppointmentResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
