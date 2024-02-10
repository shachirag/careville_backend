package doctorProfession

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetDoctorProfessionPaginationRes struct {
	Status  bool                               `json:"status"`
	Message string                             `json:"message"`
	Data    DoctorProfessionPaginationResponse `json:"data"`
}

type DoctorProfessionPaginationResponse struct {
	Total               int                      `json:"total"`
	PerPage             int                      `json:"perPage"`
	CurrentPage         int                      `json:"currentPage"`
	TotalPages          int                      `json:"totalPages"`
	DoctorProfessionRes []GetDoctorProfessionRes `json:"doctorProfessionRes"`
}

type GetDoctorProfessionRes struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Email       string             `json:"email" bson:"email"`
	FirstName   string             `json:"firstName" bson:"firstName"`
	LastName    string             `json:"lastName" bson:"lastName"`
	PhoneNumber PhoneNumber        `json:"phoneNumber" bson:"phoneNumber"`
	Speciality  string             `json:"speciality" bson:"speciality"`
	ProfileId   string             `json:"profileId" bson:"profileId"`
}

type PhoneNumber struct {
	DialCode string `json:"dialCode" bson:"dialCode"`
	Number   string `json:"number" bson:"number"`
}
