package services

import "go.mongodb.org/mongo-driver/bson/primitive"

type SubscriptionResDto struct {
	Status  bool              `json:"status" bson:"status"`
	Message string            `json:"message" bson:"message"`
	Data    []SubscriptionRes `json:"data" bson:"data"`
}

type SubscriptionRes struct {
	Id      primitive.ObjectID `json:"id" bson:"id"`
	Type    string             `json:"type" bson:"type"`
	Details string             `json:"details" bson:"details"`
	Price   float64            `json:"price" bson:"price"`
}

type GetSubscriptionResDto struct {
	Status  bool            `json:"status" bson:"status"`
	Message string          `json:"message" bson:"message"`
	Data    SubscriptionRes `json:"data" bson:"data"`
}

type SubscriptionReqDto struct {
	Type    string  `json:"type" bson:"type"`
	Details string  `json:"details" bson:"details"`
	Price   float64 `json:"price" bson:"price"`
}

type SubscriptionResponseDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}

type UpdateSubscriptionResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}

type UpdateSubscriptionReqDto struct {
	Type        string  `json:"type" bson:"type"`
	Name        string  `json:"name" bson:"name"`
	Information string  `json:"information" bson:"information"`
	Price       float64 `json:"price" bson:"price"`
}
