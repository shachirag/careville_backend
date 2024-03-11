package doctorProfession

import (
	"careville_backend/database"
	doctorProfession "careville_backend/dto/admin/services/doctorProfession"
	"careville_backend/entity"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary get-doctor detail
// @Description get-doctor detail
// @Tags admin hospitals
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} doctorProfession.GetDoctorProfessionDetailResDto
// @Router /admin/healthProfessional/get-doctor/{id} [get]
func GetDoctorDetail(c *fiber.Ctx) error {
	ctx := context.TODO()

	var service entity.ServiceEntity

	serviceColl := database.GetCollection("service")

	idParam := c.Params("id")
	hospitalID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(doctorProfession.GetDoctorProfessionDetailResDto{
			Status:  false,
			Message: "Invalid appointment ID",
		})
	}

	filter := bson.M{
		"_id": hospitalID,
	}

	projection := bson.M{
		"_id":                          1,
		"facilityOrProfession":         1,
		"role":                         1,
		"profileId":                    1,
		"serviceStatus":                1,
		"user.firstName":               1,
		"user.lastName":                1,
		"user.email":                   1,
		"user.phoneNumber.dialCode":    1,
		"user.phoneNumber.countryCode": 1,
		"user.phoneNumber.number":      1,
		"doctor.professionalDetailsDocs.certificate": 1,
		"doctor.professionalDetailsDocs.license":     1,
		"doctor.personalIdentificationDocs.nimc":     1,
		"doctor.personalIdentificationDocs.license":  1,
		"doctor.additionalServices.speciality":       1,
		"doctor.additionalServices.qualifications":   1,
		"doctor.information.name":                    1,
		"doctor.information.image":                   1,
		"doctor.information.additionalText":          1,
		"doctor.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
		"doctor.schedule.consultationFees": 1,
		"doctor.schedule.slots": bson.M{
			"id":        1,
			"startTime": 1,
			"endTime":   1,
			"days":      1,
		},
	}

	findOptions := options.FindOne().SetProjection(projection)

	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(doctorProfession.GetDoctorProfessionDetailResDto{
				Status:  false,
				Message: "Other service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(doctorProfession.GetDoctorProfessionDetailResDto{
			Status:  false,
			Message: "Failed to fetch other service from MongoDB: " + err.Error(),
		})
	}

	if service.Doctor == nil {
		return c.Status(fiber.StatusOK).JSON(doctorProfession.GetDoctorProfessionDetailResDto{
			Status:  false,
			Message: "Doctor information not found.",
		})
	}

	var slotsData []doctorProfession.Slots
	if service.Doctor != nil && len(service.Doctor.Schedule.Slots) > 0 {
		scheduleData := make([]doctorProfession.Slots, 0)
		for _, schedule := range service.Doctor.Schedule.Slots {
			scheduleData = append(scheduleData, doctorProfession.Slots{
				Id:        schedule.Id,
				StartTime: schedule.StartTime,
				EndTime:   schedule.EndTime,
				Days:      schedule.Days,
			})
		}
	}

	var doctorImage string
	var doctorName string
	var additionalText string
	var doctorAddress doctorProfession.Address
	var professionalLicense string
	var professionalCertificate string
	var personalNimc string
	var personalLicense string
	var speciality string
	var qualification string
	var consultationFees float64
	if service.Doctor != nil {
		doctorName = service.Doctor.Information.Name
		doctorImage = service.Doctor.Information.Image
		additionalText = service.Doctor.Information.AdditionalText
		doctorAddress = doctorProfession.Address(service.Doctor.Information.Address)
		professionalLicense = service.Doctor.ProfessionalDetailsDocs.License
		professionalCertificate = service.Doctor.ProfessionalDetailsDocs.Certificate
		personalNimc = service.Doctor.PersonalIdentificationDocs.Nimc
		personalLicense = service.Doctor.PersonalIdentificationDocs.License
		speciality = service.Doctor.AdditionalServices.Speciality
		qualification = service.Doctor.AdditionalServices.Qualifications
		consultationFees = service.Doctor.Schedule.ConsultationFees
	}

	response := doctorProfession.GetDoctorProfessionDetailResDto{
		Status:  true,
		Message: "Data fetched retrieved successfully",
		Data: doctorProfession.GetDoctorProfessionDetailRes{
			Id:                   service.Id,
			FacilityOrProfession: service.FacilityOrProfession,
			Role:                 service.Role,
			ProfileId:            service.ProfileId,
			ServiceStatus:        service.ServiceStatus,
			User: doctorProfession.User{
				FirstName: service.User.FirstName,
				LastName:  service.User.LastName,
				Email:     service.User.Email,
				PhoneNumber: doctorProfession.PhoneNumber{
					DialCode: service.User.PhoneNumber.DialCode,
					Number:   service.User.PhoneNumber.Number,
				},
			},
			ProfessionalDetails: doctorProfession.ProfessionalDetails{
				Speciality:     speciality,
				Qualifications: qualification,
			},
			Schedule: doctorProfession.Schedule{
				Slots:            slotsData,
				ConsultationFees: consultationFees,
			},
			DoctorProfessionInformation: doctorProfession.DoctorProfessionInformation{
				Name:           doctorName,
				Image:          doctorImage,
				AdditionalText: additionalText,
				Address:        doctorAddress,
			},
			ProfessionalDocuments: doctorProfession.ProfessionalDocuments{
				License:     professionalLicense,
				Certificate: professionalCertificate,
			},
			PersonalDocuments: doctorProfession.PersonalDocuments{
				Nimc:    personalNimc,
				License: personalLicense,
			},
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
