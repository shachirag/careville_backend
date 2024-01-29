package adminAuth

type GetAdminResDto struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    GetAdminRes `json:"data"`
}
