package schedular

type NotificationRes struct {
		Status  bool   `json:"status" bson:"status"`
		Message string `json:"message" bson:"message"`
}