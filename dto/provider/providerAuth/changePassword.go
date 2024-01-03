package providerAuth

type ProviderChangePasswordReqDto struct {
	CurrentPassword string `json:"currentPassword" bson:"currentPassword"`
	NewPassword     string `json:"newPassword" bson:"newPassword"`
	ConfirmPassword string `json:"confirmPassword" bson:"confirmPassword"`
}

type ProviderChangePasswordResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}
