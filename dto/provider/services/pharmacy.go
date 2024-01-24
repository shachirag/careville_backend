package services

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetPharmacyOtherServicesResDto struct {
	Status  bool                      `json:"status" bson:"status"`
	Message string                    `json:"message" bson:"message"`
	Data    []PharmacyOtherServiceRes `json:"data" bson:"data"`
}

type PharmacyOtherServiceRes struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	Information string             `json:"information" bson:"information"`
}

type PharmacyOtherServiceResDto struct {
	Status  bool                    `json:"status" bson:"status"`
	Message string                  `json:"message" bson:"message"`
	Data    PharmacyOtherServiceRes `json:"data" bson:"data"`
}

type UpdatePharmacyOtherServiceReqDto struct {
	Name        string `json:"name" bson:"name"`
	Information string `json:"information" bson:"information"`
}

type UpdatePharmacyOtherServiceResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}

type AddPharmacyOtherServiceReqDto struct {
	Name        string `json:"name" bson:"name"`
	Information string `json:"information" bson:"information"`
}
