package doctorProfession

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetDoctorProfessionRes struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	Image      string             `json:"image" bson:"image"`
	Name       string             `json:"name" bson:"name"`
	Speciality string             `json:"speciality" bson:"speciality"`
	AvgRating  float64            `json:"avgRating" bson:"avgRating"`
}

type GetDoctorProfessionResponseDto struct {
	Status  bool                     `json:"status"`
	Message string                   `json:"message"`
	Data    []GetDoctorProfessionRes `json:"data"`
}
