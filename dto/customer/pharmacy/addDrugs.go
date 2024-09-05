package pharmacy

import "go.mongodb.org/mongo-driver/bson/primitive"

type PharmacyDrugsReqDto struct {
	ModeOfDelivery  string `json:"modeOfDelivery" form:"modeOfDelivery"`
	NameAndQuantity string `json:"nameAndQuantity" form:"nameAndQuantity"`
	Latitude        string `json:"latitude" form:"latitude"`
	Longitude       string `json:"longitude" form:"longitude"`
	Address         string `json:"address" form:"address"`
}

type PharmacyDrugsResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type AmountPaymentForPharmacyDrugsResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type AmountPaymentForPharmacyDrugsReqDto struct {
	Amount        float64            `json:"amount"`
	AppointmentId primitive.ObjectID `json:"appointmentId"`
}
