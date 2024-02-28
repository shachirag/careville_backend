package subEntity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UpdateServiceSubEntity struct {
	Role                 string                                     `json:"role" bson:"role"`
	FacilityOrProfession string                                     `json:"facilityOrProfession" bson:"facilityOrProfession"`
	HospClinic           *HospClinicUpdateServiceSubEntity          `json:"hospClinic" bson:"hospClinic,omitempty"`
	FitnessCenter        *FitnessCenterUpdateServiceSubEntity       `json:"fitnessCenter" bson:"fitnessCenter,omitempty"`
	Laboratory           *LaboratoryUpdateServiceSubEntity          `json:"laboratory" bson:"laboratory,omitempty"`
	Pharmacy             *PharmacyUpdateServiceSubEntity            `json:"pharmacy" bson:"pharmacy,omitempty"`
	MedicalLabScientist  *MedicalLabScientistUpdateServiceSubEntity `json:"medicalLabScientist" bson:"medicalLabScientist,omitempty"`
	Doctor               *DoctorProfessionUpdateServiceSubEntity    `json:"doctor" bson:"doctor,omitempty"`
	Physiotherapist      *PhysiotherapistUpdateServiceSubEntity     `json:"physiotherapist" bson:"physiotherapist,omitempty"`
	Nurse                *NurseUpdateServiceSubEntity               `json:"nurse" bson:"nurse,omitempty"`
	ServiceStatus        string                                     `json:"serviceStatus" bson:"serviceStatus"`
	UpdatedAt            time.Time                                  `json:"updatedAt" bson:"updatedAt"`
}

type ProviderUserUpdateServiceSubEntity struct {
	FirstName    string                             `json:"firstName" bson:"firstName"`
	LastName     string                             `json:"lastName" bson:"lastName"`
	Email        string                             `json:"email" bson:"email"`
	Password     string                             `json:"password" bson:"password"`
	Notification NotificationUpdateServiceSubEntity `json:"notification" bson:"notification"`
	PhoneNumber  PhoneNumberUpdateServiceSubEntity  `json:"phoneNumber" bson:"phoneNumber"`
}

type HospClinicUpdateServiceSubEntity struct {
	Information   InformationUpdateServiceSubEntity `json:"information" bson:"information"`
	Doctor        []DoctorUpdateServiceSubEntity    `json:"doctor" bson:"doctor"`
	OtherServices []string                          `json:"otherServices" bson:"otherServices"`
	Insurances    []string                          `json:"insurances" bson:"insurances"`
	Documents     DocumentsUpdateServiceSubEntity   `json:"documents" bson:"documents"`
}

type InformationUpdateServiceSubEntity struct {
	Name                 string                        `json:"name" bson:"name"`
	AdditionalText       string                        `json:"additionalText" bson:"additionalText"`
	Image                string                        `json:"image" bson:"image"`
	Address              AddressUpdateServiceSubEntity `json:"address" bson:"address"`
	IsEmergencyAvailable bool                          `json:"isEmergencyAvailable" bson:"isEmergencyAvailable"`
}

type NotificationUpdateServiceSubEntity struct {
	DeviceToken string `json:"deviceToken" bson:"deviceToken"`
	DeviceType  string `json:"deviceType" bson:"deviceType"`
	IsEnabled   bool   `json:"isEnabled" bson:"isEnabled"`
}

type PhoneNumberUpdateServiceSubEntity struct {
	DialCode    string `json:"dialCode" bson:"dialCode"`
	CountryCode string `json:"countryCode" bson:"countryCode"`
	Number      string `json:"number" bson:"number"`
}

type AddressUpdateServiceSubEntity struct {
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
	Add         string    `json:"add" bson:"add"`
	Type        string    `json:"type" bson:"type"`
}

type DoctorUpdateServiceSubEntity struct {
	Id         primitive.ObjectID               `json:"id" bson:"id"`
	Image      string                           `json:"image" bson:"image"`
	Name       string                           `json:"name" bson:"name"`
	Speciality string                           `json:"speciality" bson:"speciality"`
	Schedule   []ScheduleUpdateServiceSubEntity `json:"schedule" bson:"schedule"`
}

type ScheduleUpdateServiceSubEntity struct {
	StartTime     string          `json:"startTime" bson:"startTime"`
	EndTime       string          `json:"endTime" bson:"endTime"`
	Days          []string        `json:"days" bson:"days"`
	BreakingSlots []BreakingSlots `json:"breakingSlots" bson:"breakingSlots"`
}

type BreakingSlots struct {
	StartTime string `json:"startTime" bson:"startTime"`
	EndTime   string `json:"endTime" bson:"endTime"`
}

type DocumentsUpdateServiceSubEntity struct {
	Certificate string `json:"certificate" bson:"certificate"`
	License     string `json:"license" bson:"license"`
}

type FitnessCenterUpdateServiceSubEntity struct {
	Information        InformationUpdateServiceSubEntity          `json:"information" bson:"information"`
	Trainers           []TrainersUpdateServiceSubEntity           `json:"trainers" bson:"trainers"`
	AdditionalServices []AdditionalServicesUpdateServiceSubEntity `json:"additionalServices" bson:"additionalServices"`
	Documents          DocumentsUpdateServiceSubEntity            `json:"documents" bson:"documents"`
	Subscription       []SubscriptionUpdateServiceSubEntity       `json:"subscription" bson:"subscription"`
}

type TrainersUpdateServiceSubEntity struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Category    string             `json:"category" bson:"category"`
	Name        string             `json:"name" bson:"name"`
	Information string             `json:"information" bson:"information"`
	Price       float64            `json:"price" bson:"price"`
}

type SubscriptionUpdateServiceSubEntity struct {
	Id      primitive.ObjectID `json:"id" bson:"id"`
	Type    string             `json:"type" bson:"type"`
	Details string             `json:"details" bson:"details"`
	Price   float64            `json:"price" bson:"price"`
}

type AdditionalServicesUpdateServiceSubEntity struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	Information string             `json:"information" bson:"information"`
}

type LaboratoryUpdateServiceSubEntity struct {
	Information    InformationUpdateServiceSubEntity      `json:"information" bson:"information"`
	Investigations []InvestigationsUpdateServiceSubEntity `json:"investigations" bson:"investigations"`
	Documents      DocumentsUpdateServiceSubEntity        `json:"documents" bson:"documents"`
}

type InvestigationsUpdateServiceSubEntity struct {
	Id          primitive.ObjectID `json:"id" bson:"id"`
	Type        string             `json:"type" bson:"type"`
	Name        string             `json:"name" bson:"name"`
	Information string             `json:"information" bson:"information"`
	Price       float64            `json:"price" bson:"price"`
}

type PharmacyUpdateServiceSubEntity struct {
	Information        InformationUpdateServiceSubEntity          `json:"information" bson:"information"`
	AdditionalServices []AdditionalServicesUpdateServiceSubEntity `json:"additionalServices" bson:"additionalServices"`
	Documents          DocumentsUpdateServiceSubEntity            `json:"documents" bson:"documents"`
}

type DoctorProfessionUpdateServiceSubEntity struct {
	Information                InformationUpdateServiceSubEntity                `json:"information" bson:"information"`
	AdditionalServices         AdditionalServiceUpdateServiceSubEntity          `json:"additionalServices" bson:"additionalServices"`
	PersonalIdentificationDocs PersonalIdentificationDocsUpdateServiceSubEntity `json:"personalIdentificationDocs" bson:"personalIdentificationDocs"`
	ProfessionalDetailsDocs    ProfessionalDetailsDocsUpdateServiceSubEntity    `json:"professionalDetailsDocs" bson:"professionalDetailsDocs"`
	Schedule                   DoctorScheduleUpdateServiceSubEntity             `json:"schedule" bson:"schedule"`
}

type AdditionalServiceUpdateServiceSubEntity struct {
	Speciality     string `json:"speciality" bson:"speciality"`
	Qualifications string `json:"qualifications" bson:"qualifications"`
}

type PersonalIdentificationDocsUpdateServiceSubEntity struct {
	Nimc    string `json:"nimc" bson:"nimc"`
	License string `json:"license" bson:"license"`
}

type ProfessionalDetailsDocsUpdateServiceSubEntity struct {
	Certificate string `json:"certificate" bson:"certificate"`
	License     string `json:"license" bson:"license"`
}

type DoctorScheduleUpdateServiceSubEntity struct {
	ConsultationFees float64                             `json:"consultationFees" bson:"consultationFees"`
	Slots            []DoctorSlotsUpdateServiceSubEntity `json:"slots" bson:"slots"`
}

type SlotsUpdateServiceSubEntity struct {
	StartTime     string          `json:"startTime" bson:"startTime"`
	EndTime       string          `json:"endTime" bson:"endTime"`
	Days          []string        `json:"days" bson:"days"`
	BreakingSlots []BreakingSlots `json:"breakingSlots" bson:"breakingSlots"`
}

type DoctorSlotsUpdateServiceSubEntity struct {
	Id            primitive.ObjectID `json:"id" bson:"id"`
	StartTime     string             `json:"startTime" bson:"startTime"`
	EndTime       string             `json:"endTime" bson:"endTime"`
	Days          []string           `json:"days" bson:"days"`
	BreakingSlots []BreakingSlots    `json:"breakingSlots" bson:"breakingSlots"`
}

type MedicalLabScientistUpdateServiceSubEntity struct {
	Information                InformationUpdateServiceSubEntity                `json:"information" bson:"information"`
	ProfessionalDetails        ProfessionalDetailUpdateServiceSubEntity         `json:"professionalDetails" bson:"professionalDetails"`
	PersonalIdentificationDocs PersonalIdentificationDocsUpdateServiceSubEntity `json:"personalIdentificationDocs" bson:"personalIdentificationDocs"`
	ProfessionalDetailsDocs    ProfessionalDetailsDocsUpdateServiceSubEntity    `json:"professionalDetailsDocs" bson:"professionalDetailsDocs"`
	ServiceAndSchedule         []ServiceAndScheduleUpdateServiceSubEntity       `json:"serviceAndSchedule" bson:"serviceAndSchedule"`
}

type ProfessionalDetailUpdateServiceSubEntity struct {
	Department    string `json:"department" bson:"department"`
	Qualification string `json:"qualification" bson:"qualification"`
}

type NurseUpdateServiceSubEntity struct {
	Information                InformationUpdateServiceSubEntity                `json:"information" bson:"information"`
	ProfessionalDetails        ProfessionalDetailsUpdateServiceSubEntity        `json:"professionalDetails" bson:"professionalDetails"`
	PersonalIdentificationDocs PersonalIdentificationDocsUpdateServiceSubEntity `json:"personalIdentificationDocs" bson:"personalIdentificationDocs"`
	ProfessionalDetailsDocs    ProfessionalDetailsDocsUpdateServiceSubEntity    `json:"professionalDetailsDocs" bson:"professionalDetailsDocs"`
	Schedule                   []ServiceAndScheduleUpdateServiceSubEntity       `json:"schedule" bson:"schedule"`
}

type ProfessionalDetailsUpdateServiceSubEntity struct {
	Qualifications string `json:"qualifications" bson:"qualifications"`
}

type ServiceAndScheduleUpdateServiceSubEntity struct {
	Id          primitive.ObjectID            `json:"id" bson:"id"`
	Name        string                        `json:"name" bson:"name"`
	ServiceFees float64                       `json:"serviceFees" bson:"serviceFees"`
	Slots       []SlotsUpdateServiceSubEntity `json:"slots" bson:"slots"`
}

type PhysiotherapistUpdateServiceSubEntity struct {
	Information                InformationUpdateServiceSubEntity                `json:"information" bson:"information"`
	ProfessionalDetails        ProfessionalDetailsUpdateServiceSubEntity        `json:"professionalDetails" bson:"professionalDetails"`
	PersonalIdentificationDocs PersonalIdentificationDocsUpdateServiceSubEntity `json:"personalIdentificationDocs" bson:"personalIdentificationDocs"`
	ProfessionalDetailsDocs    ProfessionalDetailsDocsUpdateServiceSubEntity    `json:"professionalDetailsDocs" bson:"professionalDetailsDocs"`
	ServiceAndSchedule         []ServiceAndScheduleUpdateServiceSubEntity       `json:"serviceAndSchedule" bson:"serviceAndSchedule"`
}
