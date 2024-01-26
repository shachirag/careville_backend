package providerAuth

type GetProviderResDto struct {
	Status                bool                  `json:"status"`
	Message               string                `json:"message"`
	Provider              ProviderRespDto       `json:"data"`
	AdditionalInformation AdditionalInformation `json:"additionalInformation"`
}

type AdditionalInformation struct {
	AdditionalDetails    string    `json:"additionalDetails" bson:"additionalDetails"`
	Address              Address   `json:"address" bson:"address"`
	IsEmergencyAvailable bool      `json:"isEmergencyAvailable" bson:"isEmergencyAvailable"`
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
