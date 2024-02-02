package fitnessCenter

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetFitnessCenterPaginationRes struct {
	Status  bool                            `json:"status"`
	Message string                          `json:"message"`
	Data    FitnessCenterPaginationResponse `json:"data"`
}

type FitnessCenterPaginationResponse struct {
	Total            int                   `json:"total"`
	PerPage          int                   `json:"perPage"`
	CurrentPage      int                   `json:"currentPage"`
	TotalPages       int                   `json:"totalPages"`
	FitnessCenterRes []GetFitnessCenterRes `json:"fitnessCenterRes"`
}

type GetFitnessCenterRes struct {
	Id      primitive.ObjectID `json:"id" bson:"_id"`
	Image   string             `json:"image" bson:"image"`
	Name    string             `json:"name" bson:"name"`
	Address Address            `json:"address" bson:"address"`
}

type Address struct {
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
	Add         string    `json:"add" bson:"add"`
	Type        string    `json:"type" bson:"type"`
}
