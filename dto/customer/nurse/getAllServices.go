package nurse

import "go.mongodb.org/mongo-driver/bson/primitive"

type NurseServicesResDto struct {
	Status  bool              `json:"status" bson:"status"`
	Message string            `json:"message" bson:"message"`
	Data    []NurseServiceRes `json:"data" bson:"data"`
}

type NurseServiceRes struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	ServiceFees float64            `json:"serviceFees" bson:"serviceFees"`
}
