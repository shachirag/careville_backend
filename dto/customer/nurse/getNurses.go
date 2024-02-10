package nurse

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetNurseRes struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Image     string             `json:"image" bson:"image"`
	Name      string             `json:"name" bson:"name"`
	AvgRating float64            `json:"avgRating" bson:"avgRating"`
}

type GetNurseResponseDto struct {
	Status  bool          `json:"status"`
	Message string        `json:"message"`
	Data    []GetNurseRes `json:"data"`
}
