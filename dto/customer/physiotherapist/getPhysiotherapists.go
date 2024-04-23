package physiotherapist

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetPhysiotherapistRes struct {
	Id            primitive.ObjectID `json:"id" bson:"_id"`
	Image         string             `json:"image" bson:"image"`
	Name          string             `json:"name" bson:"name"`
	ServiceType   string             `json:"serviceType" bson:"serviceType"`
	NextAvailable NextAvailable      `json:"nextAvailable" bson:"nextAvailable"`
}

type NextAvailable struct {
	StartTime string `json:"startTime" bson:"startTime"`
	// LastTime  string `json:"lastTime" bson:"lastTime"`
}

type GetPhysiotherapistResponseDto struct {
	Status  bool                    `json:"status"`
	Message string                  `json:"message"`
	Data    []GetPhysiotherapistRes `json:"data"`
}
