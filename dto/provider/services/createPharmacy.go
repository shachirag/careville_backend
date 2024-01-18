package services

type PharmacyRequestDto struct {
	PharmacyReqDto PharmacyReqDto `json:"data" form:"data"`
}

type PharmacyReqDto struct {
	InformationFirstName string               `json:"informationFirstName" form:"informationFirstName"`
	InformationLastName  string               `json:"informationLastName" form:"informationLastName"`
	Address              string               `json:"address" form:"address"`
	Longitude            string               `json:"longitude" form:"longitude"`
	Latitude             string               `json:"latitude" form:"latitude"`
	AdditionalText       string               `json:"additionalText" form:"additionalText"`
	AdditionalServices   []AdditionalServices `json:"additionalServices" form:"additionalServices"`
}

type PharmacyResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}
