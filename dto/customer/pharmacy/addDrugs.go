package pharmacy

type PharmacyDrugsReqDto struct {
	ModeOfDelivery  string  `json:"modeOfDelivery" form:"modeOfDelivery"`
	NameAndQuantity string  `json:"nameAndQuantity" form:"nameAndQuantity"`
	PricePaid       float64 `json:"pricePaid" form:"pricePaid"`
}

type PharmacyDrugsResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
