package fitnessCenter

type FitnessCenterAppointmentReqDto struct {
	FamillyMemberId        *string `json:"familyMemberId" form:"familyMemberId"`
	TrainerId              string  `json:"trainerId" form:"trainerId"`
	FamilyType             string  `json:"familyType" form:"familyType"`
	Package                string  `json:"package" form:"package"`
	MembershipSubscription float64 `json:"membershipSubscription" form:"membershipSubscription"`
	TotalAmountPaid        float64 `json:"totalAmountPaid" form:"totalAmountPaid"`
}

type FitnessCenterAppointmentResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
