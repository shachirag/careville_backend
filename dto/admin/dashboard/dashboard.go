package dashboard

type GetAllCounts struct {
	Status  bool           `json:"status"`
	Message string         `json:"message"`
	Data    DashboardCount `json:"data"`
}

type DashboardCount struct {
	HealthProfessionals HealthProfessionals `json:"healthProfessionals"`
	HealthFacilitities  HealthFacility      `json:"healthFacilitities"`
}

type HealthFacility struct {
	HospitalCount      int64 `json:"hospitalCount"`
	LaboratoryCount    int64 `json:"laboratoryCount"`
	FitnessCenterCount int64 `json:"fitnessCenterCount"`
	PharmacyCount      int64 `json:"pharmacyCount"`
}

type HealthProfessionals struct {
	DoctorCount              int64 `json:"doctorCount"`
	MedicalLabScientistCount int64 `json:"medicalLabScientistCount"`
	PhysiotherapistCount     int64 `json:"physiotherapistCount"`
	NurseCount               int64 `json:"nurseCount"`
}
