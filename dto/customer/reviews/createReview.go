package reviews

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	ReviewsReqDto struct {
		CustomerId  primitive.ObjectID `json:"customerId" bson:"customerId"`
		Description string             `json:"description" bson:"description"`
		Rating      float64            `json:"rating" bson:"rating"`
		ServiceId   primitive.ObjectID `json:"serviceId" bson:"serviceId"`
	}
	ReviewsResDto struct {
		Status  bool   `json:"status"`
		Message string `json:"message"`
	}
)
