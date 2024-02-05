package pharmacy

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetPharmacyResDto struct {
	Status  bool             `json:"status"`
	Message string           `json:"message"`
	Data    PharmacyResponse `json:"data"`
}

type PharmacyResponse struct {
	Id                 primitive.ObjectID   `json:"id" bson:"_id"`
	Image              string               `json:"image" bson:"image"`
	Name               string               `json:"name" bson:"name"`
	Address            Address              `json:"address" bson:"address"`
	AboutUs            string               `json:"aboutUs" bson:"aboutUs"`
	AdditionalServices []AdditionalServices `json:"additionalServices" bson:"additionalServices"`
	TotalReviews       int32                `json:"totalReviews" bson:"totalReviews"`
	AvgRating          float64              `json:"avgRating" bson:"avgRating"`
}

type AdditionalServices struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	Information string             `json:"information" bson:"information"`
}
