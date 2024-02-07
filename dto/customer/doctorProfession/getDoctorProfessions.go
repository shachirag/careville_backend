package doctorProfession

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetDoctorProfessionPaginationRes struct {
	Status  bool                               `json:"status"`
	Message string                             `json:"message"`
	Data    DoctorProfessionPaginationResponse `json:"data"`
}

type DoctorProfessionPaginationResponse struct {
	Total               int                      `json:"total"`
	PerPage             int                      `json:"perPage"`
	CurrentPage         int                      `json:"currentPage"`
	TotalPages          int                      `json:"totalPages"`
	DoctorProfessionRes []GetDoctorProfessionRes `json:"doctorProfessionRes"`
}

type GetDoctorProfessionRes struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	Image      string             `json:"image" bson:"image"`
	Name       string             `json:"name" bson:"name"`
	Speciality string             `json:"speciality" bson:"speciality"`
	AvgRating  float64            `json:"avgRating" bson:"avgRating"`
}

type GetDoctorProfessionResponseDto struct {
	Status  bool                     `json:"status"`
	Message string                   `json:"message"`
	Data    []GetDoctorProfessionRes `json:"data"`
}
