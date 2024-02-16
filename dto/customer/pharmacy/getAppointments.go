package pharmacy

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetPharmacyAppointmentsPaginationRes struct {
	Status  bool                                   `json:"status"`
	Message string                                 `json:"message"`
	Data    PharmacyAppointmentsPaginationResponse `json:"data"`
}

type PharmacyAppointmentsPaginationResponse struct {
	Total       int                          `json:"total"`
	PerPage     int                          `json:"perPage"`
	CurrentPage int                          `json:"currentPage"`
	TotalPages  int                          `json:"totalPages"`
	Drugs       []GetPharmacyAppointmentsRes `json:"drugs"`
}

type GetPharmacyAppointmentsRes struct {
	Id        primitive.ObjectID `json:"id" bson:"id"`
	ServiceId primitive.ObjectID `json:"serviceId" bson:"serviceId"`
	Name      string             `json:"name" bson:"name"`
	Image     string             `json:"image" bson:"image"`
	Address   Address            `json:"address" bson:"address"`
}
