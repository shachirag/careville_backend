package fitnessCenter

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetFitnessCenterResDto struct {
	Status  bool                  `json:"status"`
	Message string                `json:"message"`
	Data    FitnessCenterResponse `json:"data"`
}

type FitnessCenterResponse struct {
	Id                 primitive.ObjectID   `json:"id" bson:"_id"`
	Image              string               `json:"image" bson:"image"`
	Name               string               `json:"name" bson:"name"`
	Address            Address              `json:"address" bson:"address"`
	AboutUs            string               `json:"aboutUs" bson:"aboutUs"`
	AdditionalServices []AdditionalServices `json:"additionalServices" bson:"additionalServices"`
}

type AdditionalServices struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	Information string             `json:"information" bson:"information"`
}
