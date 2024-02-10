package physiotherapist

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetPhysiotherapistRes struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Image     string             `json:"image" bson:"image"`
	Name      string             `json:"name" bson:"name"`
	AvgRating float64            `json:"avgRating" bson:"avgRating"`
}

type GetPhysiotherapistResponseDto struct {
	Status  bool                    `json:"status"`
	Message string                  `json:"message"`
	Data    []GetPhysiotherapistRes `json:"data"`
}
