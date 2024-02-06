package fitnessCenter

type FitnessCenterAppointmentReqDto struct {
	FamillyMemberId string  `json:"familyMemberId" form:"familyMemberId"`
	TrainerId       string  `json:"trainerId" form:"trainerId"`
	FamilyType      string  `json:"familyType" form:"familyType"`
	Package         string  `json:"package" form:"package"`
	PricePaid       float64 `json:"pricePaid" form:"pricePaid"`
}

type FitnessCenterAppointmentResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
