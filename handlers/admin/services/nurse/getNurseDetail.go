package nurse

import (
	"careville_backend/database"
	nurse "careville_backend/dto/admin/services/nurse"
	"careville_backend/entity"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary get-nurse detail
// @Description get-nurse detail
// @Tags admin nurse
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} nurse.GetNurseDetailResDto
// @Router /admin/healthProfessional/get-nurse/{id} [get]
func GetNurseDetail(c *fiber.Ctx) error {
	ctx := context.TODO()

	var service entity.ServiceEntity

	serviceColl := database.GetCollection("service")

	idParam := c.Params("id")
	hospitalID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(nurse.GetNurseDetailResDto{
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
		"nurse.professionalDetailsDocs.certificate": 1,
		"nurse.professionalDetailsDocs.license":     1,
		"nurse.personalIdentificationDocs.nimc":     1,
		"nurse.personalIdentificationDocs.license":  1,
		"nurse.professionalDetails.qualifications":  1,
		"nurse.information.name":                    1,
		"nurse.information.image":                   1,
		"nurse.information.additionalText":          1,
		"nurse.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
		"nurse.schedule.id":          1,
		"nurse.schedule.name":        1,
		"nurse.schedule.serviceFees": 1,
		"nurse.schedule.slots": bson.M{
			"startTime": 1,
			"endTime":   1,
			"days":      1,
		},
	}

	findOptions := options.FindOne().SetProjection(projection)

	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(nurse.GetNurseDetailResDto{
				Status:  false,
				Message: "Other service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(nurse.GetNurseDetailResDto{
			Status:  false,
			Message: "Failed to fetch other service from MongoDB: " + err.Error(),
		})
	}

	if service.Nurse == nil {
		return c.Status(fiber.StatusOK).JSON(nurse.GetNurseDetailResDto{
			Status:  false,
			Message: "Nurse information not found.",
		})
	}

	var serviceData []nurse.ServiceAndSchedule
	if service.Nurse != nil && len(service.Nurse.Schedule) > 0 {
		for _, service := range service.Nurse.Schedule {
			scheduleData := make([]nurse.Slots, 0)
			for _, schedule := range service.Slots {
				scheduleData = append(scheduleData, nurse.Slots{
					StartTime: schedule.StartTime,
					EndTime:   schedule.EndTime,
					Days:      schedule.Days,
				})
			}
			serviceData = append(serviceData, nurse.ServiceAndSchedule{
				Id:          service.Id,
				Name:        service.Name,
				ServiceFees: service.ServiceFees,
				Slots:       scheduleData,
			})
		}
	}

	var nurseImage string
	var nurseName string
	var additionalText string
	var nurseAddress nurse.Address
	var professionalLicense string
	var professionalCertificate string
	var personalNimc string
	var personalLicense string
	var qualification string
	if service.Nurse != nil {
		nurseName = service.Nurse.Information.Name
		nurseImage = service.Nurse.Information.Image
		additionalText = service.Nurse.Information.AdditionalText
		nurseAddress = nurse.Address(service.Nurse.Information.Address)
		professionalLicense = service.Nurse.ProfessionalDetailsDocs.License
		professionalCertificate = service.Nurse.ProfessionalDetailsDocs.Certificate
		personalNimc = service.Nurse.PersonalIdentificationDocs.Nimc
		personalLicense = service.Nurse.PersonalIdentificationDocs.License
		qualification = service.Nurse.ProfessionalDetails.Qualifications
	}

	response := nurse.GetNurseDetailResDto{
		Status:  true,
		Message: "Data fetched retrieved successfully",
		Data: nurse.GetNurseDetailRes{
			Id:                   service.Id,
			FacilityOrProfession: service.FacilityOrProfession,
			Role:                 service.Role,
			ProfileId:            service.ProfileId,
			ServiceStatus:        service.ServiceStatus,
			User: nurse.User{
				FirstName: service.User.FirstName,
				LastName:  service.User.LastName,
				Email:     service.User.Email,
				PhoneNumber: nurse.PhoneNumber{
					DialCode: service.User.PhoneNumber.DialCode,
					Number:   service.User.PhoneNumber.Number,
				},
			},
			ProfessionalDetails: nurse.ProfessionalDetails{
				Qualification: qualification,
			},
			ServiceAndSchedule: serviceData,
			NurseInformation: nurse.NurseInformation{
				Name:           nurseName,
				Image:          nurseImage,
				AdditionalText: additionalText,
				Address:        nurseAddress,
			},
			ProfessionalDocuments: nurse.ProfessionalDocuments{
				License:     professionalLicense,
				Certificate: professionalCertificate,
			},
			PersonalDocuments: nurse.PersonalDocuments{
				Nimc:    personalNimc,
				License: personalLicense,
			},
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
