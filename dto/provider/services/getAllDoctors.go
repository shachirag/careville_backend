package services

type DoctorResDto struct {
	Status  bool    `json:"status" bson:"status"`
	Message string  `json:"message" bson:"message"`
	Data    []DoctorRes `json:"data" bson:"data"`
}

type  DoctorRes struct {
	
}
