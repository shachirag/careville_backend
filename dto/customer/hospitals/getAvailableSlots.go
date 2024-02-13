package hospitals

type AvailableSlotsResDto struct {
	Status  bool    `json:"status" bson:"status"`
	Message string  `json:"message" bson:"message"`
	Data    []Slots `json:"data" bson:"data"`
}

type Slots struct {
	StartTime string `json:"startTime" bson:"startTime"`
	EndTime   string `json:"endTime" bson:"endTime"`
}
