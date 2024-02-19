package pharmacy

type PharmacyDrugsReqDto struct {
	ModeOfDelivery  string `json:"modeOfDelivery" form:"modeOfDelivery"`
	NameAndQuantity string `json:"nameAndQuantity" form:"nameAndQuantity"`
	Longitude       string `json:"longitude" form:"longitude"`
	Latitude        string `json:"latitude" form:"latitude"`
	Address         string `json:"address" form:"address"`
}

type PharmacyDrugsResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
