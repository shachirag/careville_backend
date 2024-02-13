package hospitals

type HospitalClinicAppointmentReqDto struct {
	DoctorId        string  `json:"doctorId" form:"doctorId"`
	FamillyMemberId string  `json:"familyMemberId" form:"familyMemberId"`
	FromDate        string  `json:"fromDate" form:"fromDate"`
	ToDate          string  `json:"toDate" form:"toDate"`
	RemindMeBefore  string  `json:"remindMeBefore" form:"remindMeBefore"`
	AvailableTime   string  `json:"availableTime" form:"availableTime"`
	FamilyType      string  `json:"familyType" form:"familyType"`
	PricePaid       float64 `json:"pricePaid" form:"pricePaid"`
}

type HospitalClinicAppointmentResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
