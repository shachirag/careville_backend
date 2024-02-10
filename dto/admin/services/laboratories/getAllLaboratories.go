package laboratory

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetLaboratoryPaginationRes struct {
	Status  bool                         `json:"status"`
	Message string                       `json:"message"`
	Data    LaboratoryPaginationResponse `json:"data"`
}

type LaboratoryPaginationResponse struct {
	Total         int                `json:"total"`
	PerPage       int                `json:"perPage"`
	CurrentPage   int                `json:"currentPage"`
	TotalPages    int                `json:"totalPages"`
	LaboratoryRes []GetLaboratoryRes `json:"laboratoryRes"`
}

type GetLaboratoryRes struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Email       string             `json:"email" bson:"email"`
	FirstName   string             `json:"firstName" bson:"firstName"`
	LastName    string             `json:"lastName" bson:"lastName"`
	PhoneNumber PhoneNumber        `json:"phoneNumber" bson:"phoneNumber"`
	ProfileId   string             `json:"profileId" bson:"profileId"`
}

type PhoneNumber struct {
	DialCode string `json:"dialCode" bson:"dialCode"`
	Number   string `json:"number" bson:"number"`
}
