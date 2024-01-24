package services

import "go.mongodb.org/mongo-driver/bson/primitive"

type InvestigationResDto struct {
	Status  bool               `json:"status" bson:"status"`
	Message string             `json:"message" bson:"message"`
	Data    []InvestigationRes `json:"data" bson:"data"`
}

type InvestigationRes struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Type        string             `json:"type" bson:"type"`
	Name        string             `json:"name" bson:"name"`
	Information string             `json:"information" bson:"information"`
	Price       float64            `json:"price" bson:"price"`
}

type GetInvestigationResDto struct {
	Status  bool             `json:"status" bson:"status"`
	Message string           `json:"message" bson:"message"`
	Data    InvestigationRes `json:"data" bson:"data"`
}

type InvestigationReqDto struct {
	Type        string  `json:"type" bson:"type"`
	Name        string  `json:"name" bson:"name"`
	Information string  `json:"information" bson:"information"`
	Price       float64 `json:"price" bson:"price"`
}

type InvestigationResponseDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}

type UpdateInvestigationResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}

type UpdateInvestigationReqDto struct {
	Type        string  `json:"type" bson:"type"`
	Name        string  `json:"name" bson:"name"`
	Information string  `json:"information" bson:"information"`
	Price       float64 `json:"price" bson:"price"`
}


