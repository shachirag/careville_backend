package medicalLabScientist

import "go.mongodb.org/mongo-driver/bson/primitive"

type MedicalLabScientistServicesResDto struct {
	Status  bool                            `json:"status" bson:"status"`
	Message string                          `json:"message" bson:"message"`
	Data    []MedicalLabScientistServiceRes `json:"data" bson:"data"`
}

type MedicalLabScientistServiceRes struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	ServiceFees float64            `json:"serviceFees" bson:"serviceFees"`
}
