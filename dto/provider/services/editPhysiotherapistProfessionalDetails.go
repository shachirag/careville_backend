package services

import "go.mongodb.org/mongo-driver/bson/primitive"

type UpdatePhysiotherapistProfessionalInfoReqDto struct {
	Qualifications string `json:"qualifications" bson:"qualifications"`
}

type UpdatePhysiotherapistProfessionalInfoResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}

type DeletePhysiotherapistProfessionalInfoResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}

type PhysiotherapistServicesResDto struct {
	Status  bool                        `json:"status" bson:"status"`
	Message string                      `json:"message" bson:"message"`
	Data    []PhysiotherapistServiceRes `json:"data" bson:"data"`
}

type GetPhysiotherapistServicesResDto struct {
	Status  bool                      `json:"status" bson:"status"`
	Message string                    `json:"message" bson:"message"`
	Data    PhysiotherapistServiceRes `json:"data" bson:"data"`
}

type PhysiotherapistServiceRes struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	ServiceFees float64             `json:"serviceFees" bson:"serviceFees"`
	Slots       []Slots            `json:"slots" bson:"slots"`
}

type MorePhysiotherapistServiceReqDto struct {
	Name        string  `json:"name" bson:"name"`
	ServiceFees float64  `json:"serviceFees" bson:"serviceFees"`
	Slots       []Slots `json:"slots" bson:"slots"`
}

type MorePhysiotherapistServiceResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}

type UpdatePhysiotherapistServiceReqDto struct {
	Name        string  `json:"name" form:"name"`
	ServiceFees float64  `json:"serviceFees" form:"serviceFees"`
	Slots       []Slots `json:"slots" form:"slots"`
}

type UpdatePhysiotherapistServiceResDto struct {
	Status  bool   `json:"status" bson:"status"`
	Message string `json:"message" bson:"message"`
}
