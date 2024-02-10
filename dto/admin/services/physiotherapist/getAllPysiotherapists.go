package physiotherapist

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetPhysiotherapistPaginationRes struct {
	Status  bool                              `json:"status"`
	Message string                            `json:"message"`
	Data    PhysiotherapistPaginationResponse `json:"data"`
}

type PhysiotherapistPaginationResponse struct {
	Total              int                     `json:"total"`
	PerPage            int                     `json:"perPage"`
	CurrentPage        int                     `json:"currentPage"`
	TotalPages         int                     `json:"totalPages"`
	PhysiotherapistRes []GetPhysiotherapistRes `json:"physiotherapistRes"`
}

type GetPhysiotherapistRes struct {
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
