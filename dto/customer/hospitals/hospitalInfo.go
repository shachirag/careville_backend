package hospitals

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetHospitalsResDto struct {
	Status  bool              `json:"status"`
	Message string            `json:"message"`
	Data    HospitalsResponse `json:"data"`
}

type HospitalsResponse struct {
	Id                     primitive.ObjectID `json:"id" bson:"_id"`
	Image                  string             `json:"image" bson:"image"`
	Name                   string             `json:"name" bson:"name"`
	Address                Address            `json:"address" bson:"address"`
	AboutUs                string             `json:"aboutUs" bson:"aboutUs"`
	OtherServices          []string           `json:"otherServices" bson:"otherServices"`
	TotalReviews           int32              `json:"totalReviews" bson:"totalReviews"`
	AvgRating              float64            `json:"avgRating" bson:"avgRating"`
	IsCustomerFamilyMember bool               `json:"isCustomerFamilyMember" bson:"isCustomerFamilyMember"`
}
