package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceEntity struct {
	Id                   primitive.ObjectID   `json:"id" bson:"_id"`
	Role                 string               `json:"role" bson:"role"`
	FirstName            string               `json:"firstName" bson:"firstName"`
	LastName             string               `json:"lastName" bson:"lastName"`
	Email                string               `json:"email" bson:"email"`
	Password             string               `json:"password" bson:"password"`
	Notification         Notification         `json:"notification" bson:"notification"`
	PhoneNumber          PhoneNumber          `json:"phoneNumber" bson:"phoneNumber"`
	FacilityOrProfession string               `json:"facilityOrProfession" bson:"facilityOrProfession"`
	HospClinic           *HospClinic          `json:"hospClinic" bson:"hospClinic,omitempty"`
	FitnessCenter        *FitnessCenter       `json:"fitnessCenter" bson:"fitnessCenter,omitempty"`
	Laboratory           *Laboratory          `json:"laboratory" bson:"laboratory,omitempty"`
	Pharmacy             *Pharmacy            `json:"pharmacy" bson:"pharmacy,omitempty"`
	MedicalLabScientist  *MedicalLabScientist `json:"medicalLabScientist" bson:"medicalLabScientist,omitempty"`
	Doctor               *DoctorEntityDto     `json:"doctor" bson:"doctor,omitempty"`
	Physiotherapist      *Physiotherapist     `json:"physiotherapist" bson:"physiotherapist,omitempty"`
	Nurse                *Nurse               `json:"nurse" bson:"nurse,omitempty"`
	Status               string               `json:"status" bson:"status"`
	CreatedAt            time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt            time.Time            `json:"updatedAt" bson:"updatedAt"`
}

type HospClinic struct {
	Information   Information `json:"information" bson:"information"`
	Doctor        []Doctor    `json:"doctor" bson:"doctor"`
	OtherServices []string    `json:"otherServices" bson:"otherServices"`
	Insurances    []string    `json:"insurances" bson:"insurances"`
	Documents     Documents   `json:"documents" bson:"documents"`
}

type Information struct {
	FirstName            string  `json:"firstName" bson:"firstName"`
	LastName             string  `json:"lastName" bson:"lastName"`
	AdditionalText       string  `json:"additionalText" bson:"additionalText"`
	Image                string  `json:"image" bson:"image"`
	Address              Address `json:"address" bson:"address"`
	IsEmergencyAvailable bool    `json:"isEmergencyAvailable" bson:"isEmergencyAvailable"`
}

type Notification struct {
	DeviceToken string `json:"deviceToken" bson:"deviceToken"`
	DeviceType  string `json:"deviceType" bson:"deviceType"`
	IsEnabled   bool   `json:"isEnabled" bson:"isEnabled"`
}

type PhoneNumber struct {
	DialCode    string `json:"dialCode" bson:"dialCode"`
	CountryCode string `json:"countryCode" bson:"countryCode"`
	Number      string `json:"number" bson:"number"`
}

type Address struct {
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
	Add         string    `json:"add" bson:"add"`
	Type        string    `json:"type" bson:"type"`
}

type Doctor struct {
	Id         primitive.ObjectID `json:"id" bson:"id"`
	Image      string             `json:"image" bson:"image"`
	Name       string             `json:"name" bson:"name"`
	Speciality string             `json:"speciality" bson:"speciality"`
	Schedule   []Schedule         `json:"schedule" bson:"schedule"`
}

type Schedule struct {
	StartTime string   `json:"startTime" bson:"startTime"`
	EndTime   string   `json:"endTime" bson:"endTime"`
	Days      []string `json:"days" bson:"days"`
}

type Documents struct {
	Certificate string `json:"certificate" bson:"certificate"`
	License     string `json:"license" bson:"license"`
}

type FitnessCenter struct {
	Information        Information          `json:"information" bson:"information"`
	Trainers           []Trainers           `json:"trainers" bson:"trainers"`
	AdditionalServices []AdditionalServices `json:"additionalServices" bson:"additionalServices"`
	Documents          Documents            `json:"documents" bson:"documents"`
	Subscription       []Subscription       `json:"subscription" bson:"subscription"`
}

type Trainers struct {
	Category    string  `json:"category" bson:"category"`
	Name        string  `json:"name" bson:"name"`
	Information string  `json:"information" bson:"information"`
	Price       float64 `json:"price" bson:"price"`
}

type Subscription struct {
	Type    string  `json:"type" bson:"type"`
	Details string  `json:"details" bson:"details"`
	Price   float64 `json:"price" bson:"price"`
}

type AdditionalServices struct {
	Name        string `json:"name" bson:"name"`
	Information string `json:"information" bson:"information"`
}

type Laboratory struct {
	Information    Information      `json:"information" bson:"information"`
	Investigations []Investigations `json:"investigations" bson:"investigations"`
	Documents      Documents        `json:"documents" bson:"documents"`
}

type Investigations struct {
	Type        string  `json:"type" bson:"type"`
	Name        string  `json:"name" bson:"name"`
	Information string  `json:"information" bson:"information"`
	Price       float64 `json:"price" bson:"price"`
}

type Pharmacy struct {
	Information        Information          `json:"information" bson:"information"`
	AdditionalServices []AdditionalServices `json:"additionalServices" bson:"additionalServices"`
	Documents          Documents            `json:"documents" bson:"documents"`
}

type DoctorEntityDto struct {
	Information                Information                `json:"information" bson:"information"`
	AdditionalServices         AdditionalService          `json:"additionalServices" bson:"additionalServices"`
	PersonalIdentificationDocs PersonalIdentificationDocs `json:"personalIdentificationDocs" bson:"personalIdentificationDocs"`
	ProfessionalDetailsDocs    ProfessionalDetailsDocs    `json:"professionalDetailsDocs" bson:"professionalDetailsDocs"`
	Schedule                   []DoctorSchedule           `json:"schedule" bson:"schedule"`
}

type AdditionalService struct {
	Speciality     string `json:"speciality" bson:"speciality"`
	Qualifications string `json:"qualifications" bson:"qualifications"`
}

type PersonalIdentificationDocs struct {
	Nimc    string `json:"nimc" bson:"nimc"`
	License string `json:"license" bson:"license"`
}

type ProfessionalDetailsDocs struct {
	Certificate string `json:"certificate" bson:"certificate"`
	License     string `json:"license" bson:"license"`
}

type DoctorSchedule struct {
	ConsultationFees string `json:"consultationFees" bson:"consultationFees"`
	Slots            Slots  `json:"slots" bson:"slots"`
}

type Slots struct {
	StartTime string   `json:"startTime" bson:"startTime"`
	EndTime   string   `json:"endTime" bson:"endTime"`
	Days      []string `json:"days" bson:"days"`
}

type MedicalLabScientist struct {
	Information                Information                `json:"information" bson:"information"`
	PersonalDetails            PersonalDetails            `json:"personalDetails" bson:"personalDetails"`
	PersonalIdentificationDocs PersonalIdentificationDocs `json:"personalIdentificationDocs" bson:"personalIdentificationDocs"`
	ProfessionalDetailsDocs    ProfessionalDetailsDocs    `json:"professionalDetailsDocs" bson:"professionalDetailsDocs"`
	ServiceAndSchedule         []ServiceAndSchedule       `json:"serviceAndSchedule" bson:"serviceAndSchedule"`
}

type PersonalDetails struct {
	Department string `json:"department" bson:"department"`
	Document   string `json:"document" bson:"document"`
}

type Nurse struct {
	Information                Information                `json:"information" bson:"information"`
	ProfessionalDetails        ProfessionalDetails        `json:"professionalDetails" bson:"professionalDetails"`
	PersonalIdentificationDocs PersonalIdentificationDocs `json:"personalIdentificationDocs" bson:"personalIdentificationDocs"`
	ProfessionalDetailsDocs    ProfessionalDetailsDocs    `json:"professionalDetailsDocs" bson:"professionalDetailsDocs"`
	Schedule                   []ServiceAndSchedule       `json:"schedule" bson:"schedule"`
}

type ProfessionalDetails struct {
	Qualifications string `json:"qualifications" bson:"qualifications"`
}

type ServiceAndSchedule struct {
	Name        string `json:"name" bson:"name"`
	ServiceFees string `json:"serviceFees" bson:"serviceFees"`
	Slots       Slots  `json:"slots" bson:"slots"`
}

type Physiotherapist struct {
	Information                Information                `json:"information" bson:"information"`
	ProfessionalDetails        ProfessionalDetails        `json:"professionalDetails" bson:"professionalDetails"`
	PersonalIdentificationDocs PersonalIdentificationDocs `json:"personalIdentificationDocs" bson:"personalIdentificationDocs"`
	ProfessionalDetailsDocs    ProfessionalDetailsDocs    `json:"professionalDetailsDocs" bson:"professionalDetailsDocs"`
	ServiceAndSchedule         []ServiceAndSchedule       `json:"serviceAndSchedule" bson:"serviceAndSchedule"`
}
