package common

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetPastAppointmentsPaginationRes struct {
	Status  bool                                   `json:"status"`
	Message string                                 `json:"message"`
	Data    PastAppointmentsPaginationResponse `json:"data"`
}

type PastAppointmentsPaginationResponse struct {
	Total          int                      `json:"total"`
	PerPage        int                      `json:"perPage"`
	CurrentPage    int                      `json:"currentPage"`
	TotalPages     int                      `json:"totalPages"`
	AppointmentRes []GetPastAppointmentsRes `json:"pastAppointments"`
}

type GetPastAppointmentsRes struct {
	Id                   primitive.ObjectID `json:"id" bson:"id"`
	Name                 string             `json:"name" bson:"name"`
	Image                string             `json:"image" bson:"image"`
	FacilityOrProfession string             `json:"facilityOrProfession" bson:"facilityOrProfession"`
	PricePaid            float64            `json:"pricePaid" bson:"pricePaid"`
}
