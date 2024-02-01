package customerAuth

type CustomerForgotPasswordReqDto struct {
	Email string `json:"email" bson:"email"`
}

type VerifyOtpReqDto struct {
	Email      string `json:"email" bson:"email"`
	EnteredOTP string `json:"otp" bson:"otp"`
}

type ResetPasswordAfterOtpReqDto struct {
	Email           string `json:"email" bson:"email"`
	NewPassword     string `json:"newPassword" bson:"newPassword"`
	ConfirmPassword string `json:"confirmPassword" bson:"confirmPassword"`
}

type CustomerPasswordResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
