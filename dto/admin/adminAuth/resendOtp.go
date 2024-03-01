package adminAuth

type UserPasswordResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type ResendOtpReqDto struct {
	Email string `json:"email" bson:"email"`
}