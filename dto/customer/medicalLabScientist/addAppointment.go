package medicalLabScientist

type MedicalLabScientistAppointmentReqDto struct {
	FamillyMemberId string  `json:"familyMemberId" form:"familyMemberId"`
	FromDate        string  `json:"fromDate" form:"fromDate"`
	ToDate          string  `json:"toDate" form:"toDate"`
	RemindMeBefore  string  `json:"remindMeBefore" form:"remindMeBefore"`
	FamilyType      string  `json:"familyType" form:"familyType"`
	PricePaid       float64 `json:"pricePaid" form:"pricePaid"`
}

type MedicalLabScientistAppointmentResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
