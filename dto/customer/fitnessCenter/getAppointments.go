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
	Id        primitive.ObjectID `json:"id" bson:"id"`
	ServiceId primitive.ObjectID `json:"serviceId" bson:"serviceId"`
	TrainerId primitive.ObjectID `json:"trainerId" bson:"trainerId"`
	Name      string             `json:"name" bson:"name"`
	Image     string             `json:"image" bson:"image"`
	Category  string             `json:"category" bson:"category"`
}
