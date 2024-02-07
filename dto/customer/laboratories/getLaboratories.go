package laboratory

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetLaboratoryPaginationRes struct {
	Status  bool                         `json:"status"`
	Message string                       `json:"message"`
	Data    LaboratoryPaginationResponse `json:"data"`
}

type LaboratoryPaginationResponse struct {
	Total         int                `json:"total"`
	PerPage       int                `json:"perPage"`
	CurrentPage   int                `json:"currentPage"`
	TotalPages    int                `json:"totalPages"`
	LaboratoryRes []GetLaboratoryRes `json:"laboratoryRes"`
}

type GetLaboratoryRes struct {
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

type GetLaboratoryResponseDto struct {
	Status  bool               `json:"status"`
	Message string             `json:"message"`
	Data    []GetLaboratoryRes `json:"data"`
}
