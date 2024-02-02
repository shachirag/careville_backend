package physiotherapist

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetPhysiotherapistPaginationRes struct {
	Status  bool                              `json:"status"`
	Message string                            `json:"message"`
	Data    PhysiotherapistPaginationResponse `json:"data"`
}

type PhysiotherapistPaginationResponse struct {
	Total              int                     `json:"total"`
	PerPage            int                     `json:"perPage"`
	CurrentPage        int                     `json:"currentPage"`
	TotalPages         int                     `json:"totalPages"`
	PhysiotherapistRes []GetPhysiotherapistRes `json:"physiotherapistRes"`
}

type GetPhysiotherapistRes struct {
	Id    primitive.ObjectID `json:"id" bson:"_id"`
	Image string             `json:"image" bson:"image"`
	Name  string             `json:"name" bson:"name"`
}
