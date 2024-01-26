package providerAuth

type (
	ProviderSignupReqDto struct {
		Email       string `json:"email" bson:"email"`
		DialCode    string `json:"dialCode" bson:"dialCode"`
		PhoneNumber string `json:"phoneNumber" bson:"phoneNumber"`
	}
	ProviderResponseDto struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}
)

type ProviderSignupVerifyOtpReqDto struct {
	FirstName   string `json:"firstName" bson:"firstName"`
	LastName    string `json:"lastName" bson:"lastName"`
	Email       string `json:"email" bson:"email"`
	DialCode    string `json:"dialCode" bson:"dialCode"`
	PhoneNumber string `json:"phoneNumber" bson:"phoneNumber"`
	CountryCode string `json:"countryCode" bson:"countryCode"`
	DeviceToken string `json:"deviceToken" bson:"deviceToken"`
	DeviceType  string `json:"deviceType" bson:"deviceType"`
	Password    string `json:"password" bson:"password"`
	EnteredOTP  string `json:"otp" bson:"otp"`
}

type ProviderSignupVerifyOtpResDto struct {
	Status   bool                `json:"status"`
	Message  string              `json:"message"`
	Token    string              `json:"token"`
	Provider LoginProviderResDto `json:"data"`
}

type Notification struct {
	DeviceToken string `json:"deviceToken" bson:"deviceToken"`
	DeviceType  string `json:"deviceType" bson:"deviceType"`
	IsEnabled   bool   `json:"isEnabled" bson:"isEnabled"`
}

type PhoneNumber struct {
	DialCode    string `json:"dialCode" bson:"dialCode"`
	CountryCode string `json:"countryCode" bson:"countryCode"`
	Number      string `json:"number" bson:"number"`
}
