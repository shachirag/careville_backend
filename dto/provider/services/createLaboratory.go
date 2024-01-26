package services

type LaboratoryRequestDto struct {
	LaboratoryReqDto LaboratoryReqDto `json:"data" form:"data"`
}

type LaboratoryReqDto struct {
	InformationName string           `json:"informationName" form:"informationName"`
	Address         string           `json:"address" form:"address"`
	Longitude       string           `json:"longitude" form:"longitude"`
	Latitude        string           `json:"latitude" form:"latitude"`
	AdditionalText  string           `json:"additionalText" form:"additionalText"`
	Investigations  []Investigations `json:"investigations" form:"investigations"`
}

type Investigations struct {
	Type        string  `json:"type" bson:"type"`
	Name        string  `json:"name" bson:"name"`
	Information string  `json:"information" bson:"information"`
	Price       float64 `json:"price" bson:"price"`
}

type LaboratoryResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
	Role    Role   `json:"data" bson:"data"`
}
