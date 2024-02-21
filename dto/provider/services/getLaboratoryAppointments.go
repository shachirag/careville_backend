package services

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetLaboratoryAppointmentsPaginationRes struct {
	Status  bool                                     `json:"status"`
	Message string                                   `json:"message"`
	Data    LaboratoryAppointmentsPaginationResponse `json:"data"`
}

type LaboratoryAppointmentsPaginationResponse struct {
	Total          int                            `json:"total"`
	PerPage        int                            `json:"perPage"`
	CurrentPage    int                            `json:"currentPage"`
	TotalPages     int                            `json:"totalPages"`
	AppointmentRes []GetLaboratoryAppointmentsRes `json:"appointments"`
}

type GetLaboratoryAppointmentsRes struct {
	Id                   primitive.ObjectID `json:"id" bson:"id"`
	CustomrId            primitive.ObjectID `json:"customerId" bson:"customerId"`
	FirstName            string             `json:"firstName" bson:"firstName"`
	LastName             string             `json:"lastName" bson:"lastName"`
	FacilityOrProfession string             `json:"facilityOrProfession" bson:"facilityOrProfession"`
	Role                 string             `json:"role" bson:"role"`
	AppointmentDate      time.Time          `json:"appointmentDate" bson:"appointmentDate"`
}
