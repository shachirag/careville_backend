package providerAuth

type LoginProviderReqDto struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type LoginProviderResDto struct {
	Status   bool           `json:"status"`
	Message  string         `json:"message"`
	Provider ProviderResDto `json:"data"`
	Token    string         `json:"token"`
}

type GetProvideResDto struct {
	Status   bool           `json:"status"`
	Message  string         `json:"message"`
	Provider ProviderResDto `json:"data"`
}
