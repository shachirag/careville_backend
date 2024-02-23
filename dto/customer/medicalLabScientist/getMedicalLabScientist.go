package medicalLabScientist

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetMedicalLabScientistRes struct {
	Id            primitive.ObjectID `json:"id" bson:"_id"`
	Image         string             `json:"image" bson:"image"`
	Name          string             `json:"name" bson:"name"`
	ServiceType   string             `json:"serviceType" bson:"serviceType"`
	NextAvailable NextAvailable      `json:"nextAvailable" bson:"nextAvailable"`
}

type NextAvailable struct {
	StartTime string `json:"startTime" bson:"startTime"`
	LastTime  string `json:"lastTime" bson:"lastTime"`
}

type GetMedicalLabScientistResponseDto struct {
	Status  bool                        `json:"status"`
	Message string                      `json:"message"`
	Data    []GetMedicalLabScientistRes `json:"data"`
}
