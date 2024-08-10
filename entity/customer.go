package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerEntity struct {
	Id            primitive.ObjectID `json:"id" bson:"_id"`
	FirstName     string             `json:"firstName" bson:"firstName"`
	LastName      string             `json:"lastName" bson:"lastName"`
	Email         string             `json:"email" bson:"email"`
	Image         string             `json:"image" bson:"image"`
	PhoneNumber   PhoneNumber        `json:"phoneNumber" bson:"phoneNumber"`
	Address       Address            `json:"address" bson:"address"`
	Notification  Notification       `json:"notification" bson:"notification"`
	FamilyMembers []FamilyMembers    `json:"familyMembers" bson:"familyMembers"`
	Password      string             `json:"password" bson:"password"`
	Sex           string             `json:"sex" bson:"sex"`
	Age           string             `json:"age" bson:"age"`
	Wallet        Wallet             `json:"wallet" bson:"wallet"`
	IsDeleted     bool               `json:"isDeleted" bson:"isDeleted"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type Wallet struct {
	Amount string `json:"amount" bson:"amount"`
}

type FamilyMembers struct {
	Id           primitive.ObjectID `json:"id" bson:"id"`
	Name         string             `json:"name" bson:"name"`
	Age          string             `json:"age" bson:"age"`
	Sex          string             `json:"sex" bson:"sex"`
	RelationShip string             `json:"relationShip" bson:"relationShip"`
}
