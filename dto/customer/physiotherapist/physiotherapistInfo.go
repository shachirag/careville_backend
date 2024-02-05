package physiotherapist

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetPhysiotherapistResDto struct {
	Status  bool                    `json:"status"`
	Message string                  `json:"message"`
	Data    PhysiotherapistResponse `json:"data"`
}

type PhysiotherapistResponse struct {
	Id                 primitive.ObjectID   `json:"id" bson:"_id"`
	Image              string               `json:"image" bson:"image"`
	Name               string               `json:"name" bson:"name"`
	AboutMe            string               `json:"aboutMe" bson:"aboutMe"`
	ServiceAndSchedule []ServiceAndSchedule `json:"serviceAndSchedule" bson:"serviceAndSchedule"`
	TotalReviews       int32                `json:"totalReviews" bson:"totalReviews"`
	AvgRating          float64              `json:"avgRating" bson:"avgRating"`
}

type ServiceAndSchedule struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	ServiceFees float64            `json:"serviceFees" bson:"serviceFees"`
	Slots       []Slots            `json:"slots" bson:"slots"`
}

type Slots struct {
	StartTime string   `json:"startTime" bson:"startTime"`
	EndTime   string   `json:"endTime" bson:"endTime"`
	Days      []string `json:"days" bson:"days"`
}
