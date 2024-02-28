package nurse

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AvailableSlotsResDto struct {
	Status  bool              `json:"status" bson:"status"`
	Message string            `json:"message" bson:"message"`
	Data    AvailableSlotsRes `json:"data" bson:"data"`
}

type AvailableSlotsRes struct {
	Schedule        []Schedule        `json:"schedule" bson:"schedule"`
	UpcommingEvents []UpcommingEvents `json:"upcommingEvents" bson:"upcommingEvents"`
}

type Schedule struct {
	StartTime    string         `json:"startTime" bson:"startTime"`
	EndTime      string         `json:"endTime" bson:"endTime"`
	Days         []string       `json:"days" bson:"days"`
	BreakinSlots []BreakinSlots `json:"breakinSlots" bson:"breakinSlots"`
}

type BreakinSlots struct {
	StartTime string `json:"startTime" bson:"startTime"`
	EndTime   string `json:"endTime" bson:"endTime"`
}

type UpcommingEvents struct {
	Id        primitive.ObjectID `json:"id" bson:"id"`
	StartTime time.Time          `json:"startTime" bson:"startTime"`
	EndTime   time.Time          `json:"endTime" bson:"endTime"`
}
