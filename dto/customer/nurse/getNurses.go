package nurse

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetNurseRes struct {
	Id            primitive.ObjectID `json:"id" bson:"_id"`
	Image         string             `json:"image" bson:"image"`
	Name          string             `json:"name" bson:"name"`
	NextAvailable NextAvailable      `json:"nextAvailable" bson:"nextAvailable"`
}

type NextAvailable struct {
	StartTime string `json:"startTime" bson:"startTime"`
	LastTime  string `json:"lastTime" bson:"lastTime"`
}

type GetNurseResponseDto struct {
	Status  bool          `json:"status"`
	Message string        `json:"message"`
	Data    []GetNurseRes `json:"data"`
}

type AppoiynmentResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
