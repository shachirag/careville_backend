package laboratory

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetLaboratoryResDto struct {
	Status  bool               `json:"status"`
	Message string             `json:"message"`
	Data    LaboratoryResponse `json:"data"`
}

type LaboratoryResponse struct {
	Id             primitive.ObjectID `json:"id" bson:"_id"`
	Image          string             `json:"image" bson:"image"`
	Name           string             `json:"name" bson:"name"`
	Address        Address            `json:"address" bson:"address"`
	AboutUs        string             `json:"aboutUs" bson:"aboutUs"`
	Investigations []Investigations   `json:"investigations" bson:"investigations"`
}

type Investigations struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Type        string             `json:"type" bson:"type"`
	Name        string             `json:"name" bson:"name"`
	Information string             `json:"information" bson:"information"`
	Price       float64            `json:"price" bson:"price"`
}
