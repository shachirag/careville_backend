package requests

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetRequestsPaginationRes struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Data    RequestsPaginationResponse `json:"data"`
}

type RequestsPaginationResponse struct {
	Total       int              `json:"total"`
	PerPage     int              `json:"perPage"`
	CurrentPage int              `json:"currentPage"`
	TotalPages  int              `json:"totalPages"`
	Requests    []GetRequestsRes `json:"requests"`
}

type GetRequestsRes struct {
	Id                   primitive.ObjectID `json:"id" bson:"_id"`
	Role                 string             `json:"role" bson:"role"`
	FacilityOrProfession string             `json:"facilityOrProfession" bson:"facilityOrProfession"`
	FirstName            string             `json:"firstName" bson:"firstName"`
	LastName             string             `json:"lastName" bson:"lastName"`
	PhoneNumber          PhoneNumber        `json:"phoneNumber" bson:"phoneNumber"`
	ProfileId            string             `json:"profileId" bson:"profileId"`
}

type PhoneNumber struct {
	DialCode    string `json:"dialCode" bson:"dialCode"`
	Number      string `json:"number" bson:"number"`
	CountryCode string `json:"countryCode" bson:"countryCode"`
}
