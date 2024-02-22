package physiotherapist

import "go.mongodb.org/mongo-driver/bson/primitive"

type PhysiotherapistServicesResDto struct {
	Status  bool                        `json:"status" bson:"status"`
	Message string                      `json:"message" bson:"message"`
	Data    []PhysiotherapistServiceRes `json:"data" bson:"data"`
}

type PhysiotherapistServiceRes struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	ServiceFees float64            `json:"serviceFees" bson:"serviceFees"`
}
