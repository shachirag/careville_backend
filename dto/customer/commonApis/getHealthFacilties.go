package common

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetHealthFacilityResDto struct {
	Status  bool                 `json:"status"`
	Message string               `json:"message"`
	Data    HealthFacilityResDto `json:"data"`
}

type HealthFacilityResDto struct {
	HospitalRes    []GetHealthFacilityRes `json:"hospitals"`
	LaboratoryRes  []GetHealthFacilityRes `json:"laboratories"`
	Pharmacy       []GetHealthFacilityRes `json:"pharmacies"`
	FitnessCenters []GetHealthFacilityRes `json:"fitnessCenters"`
}

type GetHealthFacilityRes struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Image     string             `json:"image" bson:"image"`
	Name      string             `json:"name" bson:"name"`
	AvgRating float64            `json:"avgRating" bson:"avgRating"`
}
