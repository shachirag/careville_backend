package services

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetPhysiotherapistAppointmentsPaginationRes struct {
	Status  bool                                          `json:"status"`
	Message string                                        `json:"message"`
	Data    PhysiotherapistAppointmentsPaginationResponse `json:"data"`
}

type PhysiotherapistAppointmentsPaginationResponse struct {
	Total          int                                 `json:"total"`
	PerPage        int                                 `json:"perPage"`
	CurrentPage    int                                 `json:"currentPage"`
	TotalPages     int                                 `json:"totalPages"`
	AppointmentRes []GetPhysiotherapistAppointmentsRes `json:"appointments"`
}

type GetPhysiotherapistAppointmentsRes struct {
	Id                   primitive.ObjectID `json:"id" bson:"id"`
	CustomrId            primitive.ObjectID `json:"customerId" bson:"customerId"`
	FirstName            string             `json:"firstName" bson:"firstName"`
	LastName             string             `json:"lastName" bson:"lastName"`
	FacilityOrProfession string             `json:"facilityOrProfession" bson:"facilityOrProfession"`
	Role                 string             `json:"role" bson:"role"`
	FromDate             time.Time          `json:"fromDate" bson:"fromDate"`
	ToDate               time.Time          `json:"toDate" bson:"toDate"`
}
