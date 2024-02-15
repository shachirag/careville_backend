package laboratory

import "go.mongodb.org/mongo-driver/bson/primitive"

type InvestigationResDto struct {
	Status  bool               `json:"status" bson:"status"`
	Message string             `json:"message" bson:"message"`
	Data    []InvestigationRes `json:"data" bson:"data"`
}

type InvestigationRes struct {
	Id    primitive.ObjectID `json:"id" bson:"id"`
	Name  string             `json:"name" bson:"name"`
	Price float64            `json:"price" bson:"price"`
}
