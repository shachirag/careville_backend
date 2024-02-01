package customerAuth

type (
	CustomerSignupReqDto struct {
		Email       string `json:"email" bson:"email"`
		DialCode    string `json:"dialCode" bson:"dialCode"`
		PhoneNumber string `json:"phoneNumber" bson:"phoneNumber"`
	}
	CustomerResponseDto struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}
)

type CustomerSignupVerifyOtpReqDto struct {
	FirstName     string          `json:"firstName" bson:"firstName"`
	LastName      string          `json:"lastName" bson:"lastName"`
	Email         string          `json:"email" bson:"email"`
	PhoneNumber   PhoneNumber        `json:"phoneNumber" bson:"phoneNumber"`
	DeviceToken   string          `json:"deviceToken" bson:"deviceToken"`
	DeviceType    string          `json:"deviceType" bson:"deviceType"`
	Address       string          `json:"address" bson:"address"`
	Latitude      string          `json:"latitude" bson:"latitude"`
	Longitude     string          `json:"longitude" bson:"longitude"`
	Password      string          `json:"password" bson:"password"`
	Sex           string          `json:"sex" bson:"sex"`
	Age           string          `json:"age" bson:"age"`
	CustomerType  string          `json:"customerType" bson:"customerType"`
	FamilyMembers []FamilyMembers `json:"familyMembers" bson:"familyMembers"`
	EnteredOTP    string          `json:"otp" bson:"otp"`
}
