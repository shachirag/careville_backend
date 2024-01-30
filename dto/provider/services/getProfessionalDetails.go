package services

type GetPhysiotherapistProfessionalDetailsResDto struct {
	Status  bool                      `json:"status"`
	Message string                    `json:"message"`
	Data    PhysiotherapistDetailsRes `json:"data"`
}

type PhysiotherapistDetailsRes struct {
	Qualification           string `json:"qualification"`
	ProfessionalLicense     string `json:"professionalLicense"`
	ProfessionalCertificate string `json:"professionalCertificate"`
}
