package providerAuth

type UpdateProviderReqDto struct {
	FirstName      string `json:"firstName" form:"firstName"`
	LastName       string `json:"lastName" form:"lastName"`
	DialCode       string `json:"dialCode" form:"dialCode"`
	PhoneNumber    string `json:"phoneNumber" form:"phoneNumber"`
	CountryCode    string `json:"countryCode" form:"countryCode"`
	Latitude       string `json:"latitude" form:"latitude" bson:"latitude"`
	Longitude      string `json:"longitude" form:"longitude" bson:"longitude"`
	Address        string `json:"address" form:"address"`
	AdditionalText string `json:"additionalText" form:"additionalText"`
}

type UpdateProviderResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type UpdateImageReqDto struct {
	OldImage string `json:"oldImage" form:"oldImage"`
}
