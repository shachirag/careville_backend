package physiotherapist

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
	Id        primitive.ObjectID `json:"id" bson:"id"`
	ServiceId primitive.ObjectID `json:"serviceId" bson:"serviceId"`
	Name      string             `json:"name" bson:"name"`
	Image     string             `json:"image" bson:"image"`
	FromDate  time.Time          `json:"fromDate" bson:"fromDate"`
	ToDate    time.Time          `json:"toDate" bson:"toDate"`
	// Address Address            `json:"address" bson:"address"`
}
