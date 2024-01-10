package services

type MiscResDto struct {
	Status  bool    `json:"status" bson:"status"`
	Message string  `json:"message" bson:"message"`
	Data    MiscRes `json:"data" bson:"data"`
}

type MiscRes struct {
	HospClinic    HospClinicEntity    `json:"hospClinic" bson:"hospClinic"`
	Laboratory    LaboratoryEntity    `json:"laboratory" bson:"laboratory"`
	FitnessCenter FitnessCenterEntity `json:"fitnessCenter" bson:"fitnessCenter"`
}

type HospClinicEntity struct {
	OtherServices []string `json:"otherServices" bson:"otherServices"`
	Insurances    []string `json:"insurances" bson:"insurances"`
}

type LaboratoryEntity struct {
	Investigations []string `json:"investigations" bson:"investigations"`
}

type FitnessCenterEntity struct {
	Categories []string `json:"categories" bson:"categories"`
}
