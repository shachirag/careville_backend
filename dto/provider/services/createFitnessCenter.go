package services

type FitnessCenterRequestDto struct {
	FitnessCenterReqDto FitnessCenterReqDto `json:"data" form:"data"`
}

type FitnessCenterReqDto struct {
	InformationName    string               `json:"informationName" form:"informationName"`
	Address            string               `json:"address" form:"address"`
	Longitude          string               `json:"longitude" form:"longitude"`
	Latitude           string               `json:"latitude" form:"latitude"`
	AdditionalText     string               `json:"additionalText" form:"additionalText"`
	AdditionalServices []AdditionalServices `json:"additionalServices" form:"additionalServices"`
	Trainers           []Trainers           `json:"trainers" form:"trainers"`
	Subscription       []Subscription       `json:"subscription" form:"subscription"`
}

type AdditionalServices struct {
	Name        string `json:"name" bson:"name"`
	Information string `json:"information" bson:"information"`
}

type Trainers struct {
	Category    string  `json:"category" bson:"category"`
	Name        string  `json:"name" bson:"name"`
	Information string  `json:"information" bson:"information"`
	Price       float64 `json:"price" bson:"price"`
}

type Subscription struct {
	Type    string  `json:"type" bson:"type"`
	Details string  `json:"details" bson:"details"`
	Price   float64 `json:"price" bson:"price"`
}

type FitnessCenterResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}
