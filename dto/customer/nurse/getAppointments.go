package nurse

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetNurseAppointmentsPaginationRes struct {
	Status  bool                                `json:"status"`
	Message string                              `json:"message"`
	Data    NurseAppointmentsPaginationResponse `json:"data"`
}

type NurseAppointmentsPaginationResponse struct {
	Total          int                       `json:"total"`
	PerPage        int                       `json:"perPage"`
	CurrentPage    int                       `json:"currentPage"`
	TotalPages     int                       `json:"totalPages"`
	AppointmentRes []GetNurseAppointmentsRes `json:"appointments"`
}

type GetNurseAppointmentsRes struct {
	Id        primitive.ObjectID `json:"id" bson:"id"`
	ServiceId primitive.ObjectID `json:"serviceId" bson:"serviceId"`
	Name      string             `json:"name" bson:"name"`
	Image     string             `json:"image" bson:"image"`
}
