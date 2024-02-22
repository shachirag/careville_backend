package nurse

type NurseAppointmentReqDto struct {
	FamillyMemberId *string `json:"familyMemberId" form:"familyMemberId"`
	FromDate        string  `json:"fromDate" form:"fromDate"`
	ToDate          string  `json:"toDate" form:"toDate"`
	RemindMeBefore  string  `json:"remindMeBefore" form:"remindMeBefore"`
	FamilyType      string  `json:"familyType" form:"familyType"`
	PricePaid       float64 `json:"pricePaid" form:"pricePaid"`
	NurseServiceId  string  `json:"nurseServiceId" form:"nurseServiceId"`
	Longitude       string  `json:"longitude" form:"longitude"`
	Latitude        string  `json:"latitude" form:"latitude"`
	Address         string  `json:"address" form:"address"`
}
