package providerAuth

type UpdateProviderReqDto struct {
	Name              string `json:"name" form:"name"`
	DialCode          string `json:"dialCode" form:"dialCode"`
	PhoneNumber       string `json:"phoneNumber" form:"phoneNumber"`
	Address           string `json:"address" form:"address"`
	OldProfileImage   string `json:"oldProfileImage" form:"oldProfileImage"`
	AdditionalDetails string `json:"additionalDetails" form:"additionalDetails"`
}

type UpdateProviderResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
