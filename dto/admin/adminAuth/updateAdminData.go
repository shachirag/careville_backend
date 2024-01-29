package adminAuth

type UpdateAdminReqDto struct {
	FirstName string `json:"firstName" form:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" form:"lastName" bson:"lastName"`
	OldImage  string `json:"oldImage" form:"oldImage" bson:"oldImage"`
}

type UpdateAdminResDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
