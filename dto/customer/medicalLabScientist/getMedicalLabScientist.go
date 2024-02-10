package medicalLabScientist

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetMedicalLabScientistRes struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Image     string             `json:"image" bson:"image"`
	Name      string             `json:"name" bson:"name"`
	AvgRating float64            `json:"avgRating" bson:"avgRating"`
}

type GetMedicalLabScientistResponseDto struct {
	Status  bool                        `json:"status"`
	Message string                      `json:"message"`
	Data    []GetMedicalLabScientistRes `json:"data"`
}
