package medicalLabScientist

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetMedicalLabScientistPaginationRes struct {
	Status  bool                                  `json:"status"`
	Message string                                `json:"message"`
	Data    MedicalLabScientistPaginationResponse `json:"data"`
}

type MedicalLabScientistPaginationResponse struct {
	Total                  int                         `json:"total"`
	PerPage                int                         `json:"perPage"`
	CurrentPage            int                         `json:"currentPage"`
	TotalPages             int                         `json:"totalPages"`
	MedicalLabScientistRes []GetMedicalLabScientistRes `json:"medicalLabScientistRes"`
}

type GetMedicalLabScientistRes struct {
	Id    primitive.ObjectID `json:"id" bson:"_id"`
	Image string             `json:"image" bson:"image"`
	Name  string             `json:"name" bson:"name"`
}

