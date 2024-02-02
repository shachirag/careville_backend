package nurse

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetNursePaginationRes struct {
	Status  bool                    `json:"status"`
	Message string                  `json:"message"`
	Data    NursePaginationResponse `json:"data"`
}

type NursePaginationResponse struct {
	Total       int           `json:"total"`
	PerPage     int           `json:"perPage"`
	CurrentPage int           `json:"currentPage"`
	TotalPages  int           `json:"totalPages"`
	NurseRes    []GetNurseRes `json:"nurseRes"`
}

type GetNurseRes struct {
	Id    primitive.ObjectID `json:"id" bson:"_id"`
	Image string             `json:"image" bson:"image"`
	Name  string             `json:"name" bson:"name"`
}