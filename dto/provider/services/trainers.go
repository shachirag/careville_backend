package services

import "go.mongodb.org/mongo-driver/bson/primitive"

type TrainerResDto struct {
	Status  bool                   `json:"status" bson:"status"`
	Message string                 `json:"message" bson:"message"`
	Data    []SpecialityTrainerRes `json:"data" bson:"data"`
}

type SpecialityTrainerRes struct {
	Category string       `json:"category"`
	Trainers []TrainerRes `json:"trainers"`
}

type TrainerRes struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Category    string             `json:"category" bson:"category"`
	Name        string             `json:"name" bson:"name"`
	Information string             `json:"information" bson:"information"`
	Price       float64            `json:"price" bson:"price"`
}

type GetTrainerResDto struct {
	Status  bool       `json:"status" bson:"status"`
	Message string     `json:"message" bson:"message"`
	Data    TrainerRes `json:"data" bson:"data"`
}

type TrainerReqDto struct {
	Category    string  `json:"category" bson:"category"`
	Name        string  `json:"name" bson:"name"`
	Information string  `json:"information" bson:"information"`
	Price       float64 `json:"price" bson:"price"`
}

type TrainerResponseDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}

type UpdateTrainerResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}

type UpdateTrainerReqDto struct {
	Category    string  `json:"category" bson:"category"`
	Name        string  `json:"name" bson:"name"`
	Information string  `json:"information" bson:"information"`
	Price       float64 `json:"price" bson:"price"`
}

type TrainerOtherServiceReqDto struct {
	Name        string `json:"name" bson:"name"`
	Information string `json:"information" bson:"information"`
}

type FitnessCenterOtherServiceResDto struct {
	Status  bool                   `json:"status" bson:"status"`
	Message string                 `json:"message" bson:"message"`
	Data    FitnessOtherServiceRes `json:"data" bson:"data"`
}

type FitnessOtherServiceRes struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	Information string             `json:"information" bson:"information"`
}

type UpdateOtherServiceReqDto struct {
	Name        string `json:"name" bson:"name"`
	Information string `json:"information" bson:"information"`
}

type GetFitnessOtherServicesResDto struct {
	Status  bool                     `json:"status" bson:"status"`
	Message string                   `json:"message" bson:"message"`
	Data    []FitnessOtherServiceRes `json:"data" bson:"data"`
}
