package customerAuth

import "go.mongodb.org/mongo-driver/bson/primitive"

type AddMemberReqDto struct {
	Sex          string `json:"sex" bson:"sex"`
	Name         string `json:"name" bson:"name"`
	Age          string `json:"age" bson:"age"`
	RelationShip string `json:"relationShip" bson:"relationShip"`
}

type AddMemberResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}

type MembeResDto struct {
	Status  bool       `json:"status" bson:"status"`
	Message string     `json:"message" bson:"message"`
	Data    []MembeRes `json:"data" bson:"data"`
}

type MembeRes struct {
	Id           primitive.ObjectID `json:"id" bson:"id"`
	Sex          string             `json:"sex" bson:"sex"`
	Name         string             `json:"name" bson:"name"`
	Age          string             `json:"age" bson:"age"`
	RelationShip string             `json:"relationShip" bson:"relationShip"`
}
