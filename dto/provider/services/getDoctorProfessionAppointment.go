package services

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetDoctorAppointmentsPaginationRes struct {
	Status  bool                                 `json:"status"`
	Message string                               `json:"message"`
	Data    DoctorAppointmentsPaginationResponse `json:"data"`
}

type DoctorAppointmentsPaginationResponse struct {
	Total          int                        `json:"total"`
	PerPage        int                        `json:"perPage"`
	CurrentPage    int                        `json:"currentPage"`
	TotalPages     int                        `json:"totalPages"`
	AppointmentRes []GetDoctorAppointmentsRes `json:"appointments"`
}

type GetDoctorAppointmentsRes struct {
	Id                   primitive.ObjectID `json:"id" bson:"id"`
	CustomrId            primitive.ObjectID `json:"customerId" bson:"customerId"`
	FirstName            string             `json:"firstName" bson:"firstName"`
	LastName             string             `json:"lastName" bson:"lastName"`
	FacilityOrProfession string             `json:"facilityOrProfession" bson:"facilityOrProfession"`
	Role                 string             `json:"role" bson:"role"`
	FromDate             time.Time          `json:"fromDate" bson:"fromDate"`
	ToDate               time.Time          `json:"toDate" bson:"toDate"`
}
