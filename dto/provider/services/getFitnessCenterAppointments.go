package services

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetFitnessCenterAppointmentsPaginationRes struct {
	Status  bool                                        `json:"status"`
	Message string                                      `json:"message"`
	Data    FitnessCenterAppointmentsPaginationResponse `json:"data"`
}

type FitnessCenterAppointmentsPaginationResponse struct {
	Total          int                               `json:"total"`
	PerPage        int                               `json:"perPage"`
	CurrentPage    int                               `json:"currentPage"`
	TotalPages     int                               `json:"totalPages"`
	AppointmentRes []GetFitnessCenterAppointmentsRes `json:"appointments"`
}

type GetFitnessCenterAppointmentsRes struct {
	Id                   primitive.ObjectID `json:"id" bson:"id"`
	CustomerId           primitive.ObjectID `json:"customerId" bson:"customerId"`
	FirstName            string             `json:"firstName" bson:"firstName"`
	LastName             string             `json:"lastName" bson:"lastName"`
	FacilityOrProfession string             `json:"facilityOrProfession" bson:"facilityOrProfession"`
	Role                 string             `json:"role" bson:"role"`
	CreatedAt            time.Time          `json:"createdAt" bson:"createdAt"`
}
