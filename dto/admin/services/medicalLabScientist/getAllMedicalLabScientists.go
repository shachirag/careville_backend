package medicalLabScientist

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetMedicalLabScientistPaginationRes struct {
	Status  bool                                  `json:"status"`
	Message string                                `json:"message"`
	Data    MedicalLabScientistPaginationResponse `json:"data"`
}

type MedicalLabScientistPaginationResponse struct {
	Total                  int                         `json:"total"`
	PerPage                int                         `json:"perPage"`
	CurrentPage            int                         `json:"currentPage"`
	TotalPages             int                         `json:"totalPages"`
	MedicalLabScientistRes []GetMedicalLabScientistRes `json:"medicalLabScientistRes"`
}

type GetMedicalLabScientistRes struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Email       string             `json:"email" bson:"email"`
	FirstName   string             `json:"firstName" bson:"firstName"`
	LastName    string             `json:"lastName" bson:"lastName"`
	PhoneNumber PhoneNumber        `json:"phoneNumber" bson:"phoneNumber"`
	Department  string             `json:"department" bson:"department"`
	ProfileId   string             `json:"profileId" bson:"profileId"`
}

type PhoneNumber struct {
	DialCode string `json:"dialCode" bson:"dialCode"`
	Number   string `json:"number" bson:"number"`
}
