package doctorProfession

import "go.mongodb.org/mongo-driver/bson/primitive"

type DoctorProfessionAppointmentReqDto struct {
	FamillyMemberId string  `json:"familyMemberId" form:"familyMemberId"`
	FromDate        string  `json:"fromDate" form:"fromDate"`
	ToDate          string  `json:"toDate" form:"toDate"`
	RemindMeBefore  string  `json:"remindMeBefore" form:"remindMeBefore"`
	FamilyType      string  `json:"familyType" form:"familyType"`
	PricePaid       float64 `json:"pricePaid" form:"pricePaid"`
	Address         string  `json:"address" form:"address"`
	Longitude       string  `json:"longitude" form:"longitude"`
	Latitude        string  `json:"latitude" form:"latitude"`
}

type DoctorProfessionAppointmentResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type GetDoctorAppointmentsPaginationRes struct {
	Status  bool                                 `json:"status"`
	Message string                               `json:"message"`
	Data    DoctorAppointmentsPaginationResponse `json:"data"`
}

type DoctorAppointmentsPaginationResponse struct {
	Total          int                        `json:"total"`
	PerPage        int                        `json:"perPage"`
	CurrentPage    int                        `json:"currentPage"`
	TotalPages     int                        `json:"totalPages"`
	AppointmentRes []GetDoctorAppointmentsRes `json:"appointments"`
}

type GetDoctorAppointmentsRes struct {
	Id         primitive.ObjectID `json:"id" bson:"id"`
	ServiceId  primitive.ObjectID `json:"seviceId" bson:"seviceId"`
	Name       string             `json:"name" bson:"name"`
	Image      string             `json:"image" bson:"image"`
	Speciality string             `json:"speciality" bson:"speciality"`
}
