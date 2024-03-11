package medicalLabScientist

import (
	"careville_backend/database"
	medicalLabScientist "careville_backend/dto/admin/services/medicalLabScientist"
	"careville_backend/entity"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary get-medicalLabScientist detail
// @Description get-medicalLabScientist detail
// @Tags admin medicalLabScientist
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} medicalLabScientist.GetMedicalLabScientistDetailResDto
// @Router /admin/healthProfessional/get-medicalLabScientist/{id} [get]
func GetMedicalLabScientistDetail(c *fiber.Ctx) error {
	ctx := context.TODO()

	var service entity.ServiceEntity

	serviceColl := database.GetCollection("service")

	idParam := c.Params("id")
	hospitalID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(medicalLabScientist.GetMedicalLabScientistDetailResDto{
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
		"medicalLabScientist.professionalDetailsDocs.certificate": 1,
		"medicalLabScientist.professionalDetailsDocs.license":     1,
		"medicalLabScientist.personalIdentificationDocs.nimc":     1,
		"medicalLabScientist.personalIdentificationDocs.license":  1,
		"medicalLabScientist.professionalDetails.department":      1,
		"medicalLabScientist.professionalDetails.qualification":   1,
		"medicalLabScientist.information.name":                    1,
		"medicalLabScientist.information.image":                   1,
		"medicalLabScientist.information.additionalText":          1,
		"medicalLabScientist.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
		"medicalLabScientist.serviceAndSchedule.id":          1,
		"medicalLabScientist.serviceAndSchedule.name":        1,
		"medicalLabScientist.serviceAndSchedule.serviceFees": 1,
		"medicalLabScientist.serviceAndSchedule.slots": bson.M{
			"startTime": 1,
			"endTime":   1,
			"days":      1,
		},
	}

	findOptions := options.FindOne().SetProjection(projection)

	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(medicalLabScientist.GetMedicalLabScientistDetailResDto{
				Status:  false,
				Message: "Other service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.GetMedicalLabScientistDetailResDto{
			Status:  false,
			Message: "Failed to fetch other service from MongoDB: " + err.Error(),
		})
	}

	if service.MedicalLabScientist == nil {
		return c.Status(fiber.StatusOK).JSON(medicalLabScientist.GetMedicalLabScientistDetailResDto{
			Status:  false,
			Message: "MedicalLabScientist information not found.",
		})
	}

	// Inside your GetMedicalLabScientistDetail function
	var serviceData []medicalLabScientist.ServiceAndSchedule
	if service.MedicalLabScientist != nil && len(service.MedicalLabScientist.ServiceAndSchedule) > 0 {
		for _, service := range service.MedicalLabScientist.ServiceAndSchedule {
			scheduleData := make([]medicalLabScientist.Slots, 0)
			for _, schedule := range service.Slots {
				scheduleData = append(scheduleData, medicalLabScientist.Slots{
					StartTime: schedule.StartTime,
					EndTime:   schedule.EndTime,
					Days:      schedule.Days,
				})
			}
			serviceData = append(serviceData, medicalLabScientist.ServiceAndSchedule{
				Id:          service.Id,
				Name:        service.Name,
				ServiceFees: service.ServiceFees,
				Slots:       scheduleData,
			})
		}
	}

	var medicalLabScientistImage string
	var medicalLabScientistName string
	var additionalText string
	var medicalLabScientistAddress medicalLabScientist.Address
	var professionalLicense string
	var professionalCertificate string
	var personalNimc string
	var personalLicense string
	var department string
	var qualification string
	if service.MedicalLabScientist != nil {
		medicalLabScientistName = service.MedicalLabScientist.Information.Name
		medicalLabScientistImage = service.MedicalLabScientist.Information.Image
		additionalText = service.MedicalLabScientist.Information.AdditionalText
		medicalLabScientistAddress = medicalLabScientist.Address(service.MedicalLabScientist.Information.Address)
		professionalLicense = service.MedicalLabScientist.ProfessionalDetailsDocs.License
		professionalCertificate = service.MedicalLabScientist.ProfessionalDetailsDocs.Certificate
		personalNimc = service.MedicalLabScientist.PersonalIdentificationDocs.Nimc
		personalLicense = service.MedicalLabScientist.PersonalIdentificationDocs.License
		department = service.MedicalLabScientist.ProfessionalDetails.Department
		qualification = service.MedicalLabScientist.ProfessionalDetails.Qualification
	}

	response := medicalLabScientist.GetMedicalLabScientistDetailResDto{
		Status:  true,
		Message: "Data fetched retrieved successfully",
		Data: medicalLabScientist.GetMedicalLabScientistDetailRes{
			Id:                   service.Id,
			FacilityOrProfession: service.FacilityOrProfession,
			Role:                 service.Role,
			ProfileId:            service.ProfileId,
			ServiceStatus:        service.ServiceStatus,
			User: medicalLabScientist.User{
				FirstName: service.User.FirstName,
				LastName:  service.User.LastName,
				Email:     service.User.Email,
				PhoneNumber: medicalLabScientist.PhoneNumber{
					DialCode: service.User.PhoneNumber.DialCode,
					Number:   service.User.PhoneNumber.Number,
				},
			},
			ProfessionalDetails: medicalLabScientist.ProfessionalDetails{
				Department:    department,
				Qualification: qualification,
			},
			ServiceAndSchedule: serviceData,
			MedicalLabScientistInformation: medicalLabScientist.MedicalLabScientistInformation{
				Name:           medicalLabScientistName,
				Image:          medicalLabScientistImage,
				AdditionalText: additionalText,
				Address:        medicalLabScientistAddress,
			},
			ProfessionalDocuments: medicalLabScientist.ProfessionalDocuments{
				License:     professionalLicense,
				Certificate: professionalCertificate,
			},
			PersonalDocuments: medicalLabScientist.PersonalDocuments{
				Nimc:    personalNimc,
				License: personalLicense,
			},
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
