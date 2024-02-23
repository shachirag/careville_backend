package fitnessCenter

import "go.mongodb.org/mongo-driver/bson/primitive"

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
	Id              primitive.ObjectID `json:"id" bson:"id"`
	FitnessCenterId primitive.ObjectID `json:"fitnessCenterId" bson:"fitnessCenterId"`
	Name            string             `json:"name" bson:"name"`
	Image           string             `json:"image" bson:"image"`
	Address         Address            `json:"address" bson:"address"`
}
