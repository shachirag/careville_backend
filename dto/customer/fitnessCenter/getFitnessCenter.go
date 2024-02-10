package fitnessCenter

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetFitnessCenterRes struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Image     string             `json:"image" bson:"image"`
	Name      string             `json:"name" bson:"name"`
	Address   Address            `json:"address" bson:"address"`
	AvgRating float64            `json:"avgRating" bson:"avgRating"`
}

type Address struct {
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
	Add         string    `json:"add" bson:"add"`
	Type        string    `json:"type" bson:"type"`
}

type GetFitnessCenterResponseDto struct {
	Status  bool                  `json:"status"`
	Message string                `json:"message"`
	Data    []GetFitnessCenterRes `json:"data"`
}
