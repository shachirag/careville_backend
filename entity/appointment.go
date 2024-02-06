package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentEntity struct {
	Id                   primitive.ObjectID                    `json:"id" bson:"_id"`
	Customer             CustomerAppointmentEntity             `json:"customer" bson:"customer"`
	ServiceID            primitive.ObjectID                    `json:"serviceId" bson:"serviceId"`
	Pharmacy             *PharmacyAppointmentEntity            `json:"pharmacy,omitempty" bson:"pharmacy,omitempty"`
	Laboratory           *LaboratoryAppointmentEntity          `json:"laboratory,omitempty" bson:"laboratory,omitempty"`
	FitnessCenter        *FitnessCenterAppointmentEntity       `json:"fitnessCenter,omitempty" bson:"fitnessCenter,omitempty"`
	HospitalClinic       *HospitalAppointmentEntity            `json:"hospital,omitempty" bson:"hospital,omitempty"`
	Nurse                *NurseAppointmentEntity               `json:"nurse,omitempty" bson:"nurse,omitempty"`
	Physiotherapist      *PhysiotherapistAppointmentEntity     `json:"physiotherapist,omitempty" bson:"physiotherapist,omitempty"`
	MedicalLabScientist  *MedicalLabScientistAppointmentEntity `json:"medicalLabScientist,omitempty" bson:"medicalLabScientist,omitempty"`
	Doctor               *DoctorProfessionAppointmentEntity    `json:"doctor,omitempty" bson:"doctor,omitempty"`
	FacilityOrProfession string                                `json:"facilityOrProfession" bson:"facilityOrProfession"`
	Role                 string                                `json:"role" bson:"role"`
	PaymentStatus        string                                `json:"paymentStatus" bson:"paymentStatus"`
	AppointmentStatus    string                                `json:"appointmentStatus" bson:"appointmentStatus"`
	CreatedAt            time.Time                             `json:"createdAt" bson:"createdAt"`
	UpdatedAt            time.Time                             `json:"updatedAt" bson:"updatedAt"`
}

type CustomerAppointmentEntity struct {
	ID          primitive.ObjectID `json:"id" bson:"id"`
	FirstName   string             `json:"firstName" bson:"firstName"`
	LastName    string             `json:"lastName" bson:"lastName"`
	Image       string             `json:"image" bson:"image"`
	Email       string             `json:"email" bson:"email"`
	PhoneNumber PhoneNumber        `json:"phoneNumber" bson:"phoneNumber"`
}

type PharmacyAppointmentEntity struct {
	RequestedDrugs RequestedDrugsAppointmentEntity      `json:"requestedDrugs" bson:"requestedDrugs"`
	PricePaid      float64                              `json:"pricePaid" bson:"pricePaid"`
	Information    PharmacyInformationAppointmentEntity `json:"information" bson:"information"`
}

type RequestedDrugsAppointmentEntity struct {
	ModeOfDelivery  string   `json:"modeOfDelivery" bson:"modeOfDelivery"`
	NameAndQuantity string   `json:"nameAndQuantity" bson:"nameAndQuantity"`
	Prescription    []string `json:"prescription" bson:"prescription"`
}

type PharmacyInformationAppointmentEntity struct {
	Name    string  `json:"name" bson:"name"`
	Image   string  `json:"image" bson:"image"`
	Address Address `json:"address" bson:"address"`
}

type LaboratoryAppointmentEntity struct {
	Investigation      InvestigationAppointmentEntity                `json:"investigation" bson:"investigation"`
	FamilyMember       FamilyMemberAppointmentEntity                 `json:"familyMember" bson:"familyMember"`
	AppointmentDetails LaboratoryAppointmentDetailsAppointmentEntity `json:"appointmentDetails" bson:"appointmentDetails"`
	FamilyType         string                                        `json:"familyType" bson:"familyType"`
	PricePaid          float64                                       `json:"pricePaid" bson:"pricePaid"`
}

type InvestigationAppointmentEntity struct {
	ID          primitive.ObjectID `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	Information string             `json:"information" bson:"information"`
	Type        string             `json:"type" bson:"type"`
	Price       float64            `json:"price" bson:"price"`
}

type FamilyMemberAppointmentEntity struct {
	ID           primitive.ObjectID `json:"id" bson:"id"`
	Name         string             `json:"name" bson:"name"`
	Age          string             `json:"age" bson:"age"`
	Relationship string             `json:"relationship" bson:"relationship"`
	Sex          string             `json:"sex" bson:"sex"`
}

type LaboratoryAppointmentDetailsAppointmentEntity struct {
	Date time.Time `json:"date" bson:"date"`
}

type FitnessCenterAppointmentEntity struct {
	Package      string                        `json:"package" bson:"package"`
	Trainer      TrainerAppointmentEntity      `json:"trainer" bson:"trainer"`
	FamilyType   string                        `json:"familyType" bson:"familyType"`
	FamilyMember FamilyMemberAppointmentEntity `json:"familyMember" bson:"familyMember"`
	PricePaid    float64                       `json:"pricePaid" bson:"pricePaid"`
}

type TrainerAppointmentEntity struct {
	ID          primitive.ObjectID `json:"id" bson:"id"`
	Category    string             `json:"category" bson:"category"`
	Name        string             `json:"name" bson:"name"`
	Information string             `json:"information" bson:"information"`
	Price       float64            `json:"price" bson:"price"`
}

type HospitalAppointmentEntity struct {
	Doctor             DoctorAppointmentEntity             `json:"doctor" bson:"doctor"`
	AppointmentDetails AppointmentDetailsAppointmentEntity `json:"appointmentDetails" bson:"appointmentDetails"`
	FamilyMember       FamilyMemberAppointmentEntity       `json:"familyMember" bson:"familyMember"`
	FamilyType         string                              `json:"familyType" bson:"familyType"`
	PricePaid          float64                             `json:"pricePaid" bson:"pricePaid"`
}

type DoctorAppointmentEntity struct {
	ID         primitive.ObjectID `json:"id" bson:"id"`
	Name       string             `json:"name" bson:"name"`
	Image      string             `json:"image" bson:"image"`
	Speciality string             `json:"speciality" bson:"speciality"`
}

type NurseAppointmentEntity struct {
	AppointmentDetails AppointmentDetailsAppointmentEntity `json:"appointmentDetails" bson:"appointmentDetails"`
	FamilyMember       FamilyMemberAppointmentEntity       `json:"familyMember" bson:"familyMember"`
	FamilyType         string                              `json:"familyType" bson:"familyType"`
	PricePaid          float64                             `json:"pricePaid" bson:"pricePaid"`
}

type PhysiotherapistAppointmentEntity struct {
	AppointmentDetails AppointmentDetailsAppointmentEntity `json:"appointmentDetails" bson:"appointmentDetails"`
	FamilyMember       FamilyMemberAppointmentEntity       `json:"familyMember" bson:"familyMember"`
	FamilyType         string                              `json:"familyType" bson:"familyType"`
	PricePaid          float64                             `json:"pricePaid" bson:"pricePaid"`
}

type MedicalLabScientistAppointmentEntity struct {
	AppointmentDetails AppointmentDetailsAppointmentEntity `json:"appointmentDetails" bson:"appointmentDetails"`
	FamilyMember       FamilyMemberAppointmentEntity       `json:"familyMember" bson:"familyMember"`
	FamilyType         string                              `json:"familyType" bson:"familyType"`
	PricePaid          float64                             `json:"pricePaid" bson:"pricePaid"`
}

type DoctorAppointment struct {
	AppointmentDetails AppointmentDetailsAppointmentEntity `json:"appointmentDetails" bson:"appointmentDetails"`
	FamilyMember       FamilyMemberAppointmentEntity       `json:"familyMember" bson:"familyMember"`
	FamilyType         string                              `json:"familyType" bson:"familyType"`
	PricePaid          float64                             `json:"pricePaid" bson:"pricePaid"`
}

type DoctorProfessionAppointmentEntity struct {
	FamilyMember       FamilyMemberAppointmentEntity       `json:"familyMember" bson:"familyMember"`
	AppointmentDetails AppointmentDetailsAppointmentEntity `json:"appointmentDetails" bson:"appointmentDetails"`
	FamilyType         string                              `json:"familyType" bson:"familyType"`
	PricePaid          float64                             `json:"pricePaid" bson:"pricePaid"`
	Destination        Address                             `json:"destination" bson:"destination"`
}

type AppointmentDetailsAppointmentEntity struct {
	Date           time.Time `json:"date" bson:"date"`
	RemindMeBefore string    `json:"remindMeBefore" bson:"remindMeBefore"`
	AvailableTime  string    `json:"availableTime" bson:"availableTime"`
}

// type AppointmentStatusAppointmentEntity struct {
// 	Time    time.Time `json:"time" bson:"time"`
// 	Status  string    `json:"status" bson:"status"`
// 	Message string    `json:"message" bson:"message"`
// }

// type PaymentStatusAppointmentEntity struct {
// 	Status  string    `json:"status" bson:"status"`
// 	Time    time.Time `json:"time" bson:"time"`
// 	Message string    `json:"message" bson:"message"`
// }
