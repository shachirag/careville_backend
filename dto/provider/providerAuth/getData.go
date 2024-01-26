package providerAuth

type GetProviderResDto struct {
	Status   bool           `json:"status"`
	Message  string         `json:"message"`
	Provider ProviderResDto `json:"data"`
}

type ProviderResDto struct {
	AdditionalInformation AdditionalInformation `json:"additionalInformation"`
	User                  UserData              `json:"userData"`
}
type UserData struct {
	Role Role `json:"role"`
	User User `json:"user"`
}

type AdditionalInformation struct {
	AdditionalDetails    string    `json:"additionalDetails" bson:"additionalDetails"`
	Address              Address   `json:"address" bson:"address"`
	Documents            Documents `json:"documents" bson:"documents"`
}

type Documents struct {
	Certificate string `json:"certificate" bson:"certificate"`
	License     string `json:"license" bson:"license"`
}
type Address struct {
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
	Add         string    `json:"add" bson:"add"`
	Type        string    `json:"type" bson:"type"`
}
