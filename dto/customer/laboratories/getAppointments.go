package laboratory

import "go.mongodb.org/mongo-driver/bson/primitive"

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
	Id              primitive.ObjectID `json:"id" bson:"id"`
	ServiceId       primitive.ObjectID `json:"serviceId" bson:"serviceId"`
	InvestigationID primitive.ObjectID `json:"investigationId" bson:"investigationId"`
	Name            string             `json:"name" bson:"name"`
	Image           string             `json:"image" bson:"image"`
	Type            string             `json:"type" bson:"type"`
}
