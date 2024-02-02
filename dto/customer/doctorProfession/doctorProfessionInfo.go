package doctorProfession

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetDoctorProfessionResDto struct {
	Status  bool                     `json:"status"`
	Message string                   `json:"message"`
	Data    DoctorProfessionResponse `json:"data"`
}

type DoctorProfessionResponse struct {
	Id               primitive.ObjectID `json:"id" bson:"_id"`
	Image            string             `json:"image" bson:"image"`
	Name             string             `json:"name" bson:"name"`
	Speciality       string             `json:"speciality" bson:"speciality"`
	AboutMe          string             `json:"aboutMe" bson:"aboutMe"`
	ConsultationFees float64            `json:"consultationFees" bson:"consultationFees"`
	DoctorSchedule   []DoctorSchedule   `json:"schedule" bson:"schedule"`
}

type DoctorSchedule struct {
	Id        primitive.ObjectID `json:"id" bson:"id"`
	StartTime string             `json:"startTime" bson:"startTime"`
	EndTime   string             `json:"endTime" bson:"endTime"`
	Days      []string           `json:"days" bson:"days"`
}
