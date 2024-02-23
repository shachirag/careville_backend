package medicalLabScientist

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetMedicalLabScientistAppointmentsPaginationRes struct {
	Status  bool                                              `json:"status"`
	Message string                                            `json:"message"`
	Data    MedicalLabScientistAppointmentsPaginationResponse `json:"data"`
}

type MedicalLabScientistAppointmentsPaginationResponse struct {
	Total          int                                     `json:"total"`
	PerPage        int                                     `json:"perPage"`
	CurrentPage    int                                     `json:"currentPage"`
	TotalPages     int                                     `json:"totalPages"`
	AppointmentRes []GetMedicalLabScientistAppointmentsRes `json:"appointments"`
}

type GetMedicalLabScientistAppointmentsRes struct {
	Id         primitive.ObjectID `json:"id" bson:"id"`
	ServiceId  primitive.ObjectID `json:"serviceId" bson:"serviceId"`
	Name       string             `json:"name" bson:"name"`
	Image      string             `json:"image" bson:"image"`
	Department string             `json:"department" bson:"department"`
	FromDate   time.Time          `json:"fromDate" bson:"fromDate"`
	ToDate     time.Time          `json:"toDate" bson:"toDate"`
}
