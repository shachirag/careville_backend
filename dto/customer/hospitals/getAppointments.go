package hospitals

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetHospitalAppointmentsPaginationRes struct {
	Status  bool                                   `json:"status"`
	Message string                                 `json:"message"`
	Data    HospitalAppointmentsPaginationResponse `json:"data"`
}

type HospitalAppointmentsPaginationResponse struct {
	Total          int                          `json:"total"`
	PerPage        int                          `json:"perPage"`
	CurrentPage    int                          `json:"currentPage"`
	TotalPages     int                          `json:"totalPages"`
	AppointmentRes []GetHospitalAppointmentsRes `json:"appointments"`
}

type GetHospitalAppointmentsRes struct {
	Id         primitive.ObjectID `json:"id" bson:"id"`
	HospitalId primitive.ObjectID `json:"hospitalId" bson:"hospitalId"`
	Name       string             `json:"name" bson:"name"`
	Image      string             `json:"image" bson:"image"`
	Address    Address            `json:"address" bson:"address"`
}
