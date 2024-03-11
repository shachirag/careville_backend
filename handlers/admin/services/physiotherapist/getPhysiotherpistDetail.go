package physiotherapist

import (
	"careville_backend/database"
	physiotherapist "careville_backend/dto/admin/services/physiotherapist"
	"careville_backend/entity"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary get-physiotherapist detail
// @Description get-physiotherapist detail
// @Tags admin physiotherapist
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} physiotherapist.GetPhysiotherapistDetailResDto
// @Router /admin/healthProfessional/get-physiotherapist/{id} [get]
func GetPhysiotherapistDetail(c *fiber.Ctx) error {
	ctx := context.TODO()

	var service entity.ServiceEntity

	serviceColl := database.GetCollection("service")

	idParam := c.Params("id")
	hospitalID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(physiotherapist.GetPhysiotherapistDetailResDto{
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
		"physiotherapist.professionalDetailsDocs.certificate": 1,
		"physiotherapist.professionalDetailsDocs.license":     1,
		"physiotherapist.personalIdentificationDocs.nimc":     1,
		"physiotherapist.personalIdentificationDocs.license":  1,
		"physiotherapist.professionalDetails.qualifications":  1,
		"physiotherapist.information.name":                    1,
		"physiotherapist.information.image":                   1,
		"physiotherapist.information.additionalText":          1,
		"physiotherapist.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
		"physiotherapist.serviceAndSchedule.id":          1,
		"physiotherapist.serviceAndSchedule.name":        1,
		"physiotherapist.serviceAndSchedule.serviceFees": 1,
		"physiotherapist.serviceAndSchedule.slots": bson.M{
			"startTime": 1,
			"endTime":   1,
			"days":      1,
		},
	}

	findOptions := options.FindOne().SetProjection(projection)

	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(physiotherapist.GetPhysiotherapistDetailResDto{
				Status:  false,
				Message: "Other service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(physiotherapist.GetPhysiotherapistDetailResDto{
			Status:  false,
			Message: "Failed to fetch other service from MongoDB: " + err.Error(),
		})
	}

	if service.Physiotherapist == nil {
		return c.Status(fiber.StatusOK).JSON(physiotherapist.GetPhysiotherapistDetailResDto{
			Status:  false,
			Message: "Physiotherapist information not found.",
		})
	}

	var serviceData []physiotherapist.ServiceAndSchedule
	if service.Physiotherapist != nil && len(service.Physiotherapist.ServiceAndSchedule) > 0 {
		for _, service := range service.Physiotherapist.ServiceAndSchedule {
			scheduleData := make([]physiotherapist.Slots, 0)
			for _, schedule := range service.Slots {
				scheduleData = append(scheduleData, physiotherapist.Slots{
					StartTime: schedule.StartTime,
					EndTime:   schedule.EndTime,
					Days:      schedule.Days,
				})
			}
			serviceData = append(serviceData, physiotherapist.ServiceAndSchedule{
				Id:          service.Id,
				Name:        service.Name,
				ServiceFees: service.ServiceFees,
				Slots:       scheduleData,
			})
		}
	}

	var physiotherapistImage string
	var physiotherapistName string
	var additionalText string
	var physiotherapistAddress physiotherapist.Address
	var professionalLicense string
	var professionalCertificate string
	var personalNimc string
	var personalLicense string
	var qualification string
	if service.Physiotherapist != nil {
		physiotherapistName = service.Physiotherapist.Information.Name
		physiotherapistImage = service.Physiotherapist.Information.Image
		additionalText = service.Physiotherapist.Information.AdditionalText
		physiotherapistAddress = physiotherapist.Address(service.Physiotherapist.Information.Address)
		professionalLicense = service.Physiotherapist.ProfessionalDetailsDocs.License
		professionalCertificate = service.Physiotherapist.ProfessionalDetailsDocs.Certificate
		personalNimc = service.Physiotherapist.PersonalIdentificationDocs.Nimc
		personalLicense = service.Physiotherapist.PersonalIdentificationDocs.License
		qualification = service.Physiotherapist.ProfessionalDetails.Qualifications
	}

	response := physiotherapist.GetPhysiotherapistDetailResDto{
		Status:  true,
		Message: "Data fetched retrieved successfully",
		Data: physiotherapist.GetPhysiotherapistDetailRes{
			Id:                   service.Id,
			FacilityOrProfession: service.FacilityOrProfession,
			Role:                 service.Role,
			ProfileId:            service.ProfileId,
			ServiceStatus:        service.ServiceStatus,
			User: physiotherapist.User{
				FirstName: service.User.FirstName,
				LastName:  service.User.LastName,
				Email:     service.User.Email,
				PhoneNumber: physiotherapist.PhoneNumber{
					DialCode: service.User.PhoneNumber.DialCode,
					Number:   service.User.PhoneNumber.Number,
				},
			},
			ProfessionalDetails: physiotherapist.ProfessionalDetails{
				Qualification: qualification,
			},
			ServiceAndSchedule: serviceData,
			PhysiotherapistInformation: physiotherapist.PhysiotherapistInformation{
				Name:           physiotherapistName,
				Image:          physiotherapistImage,
				AdditionalText: additionalText,
				Address:        physiotherapistAddress,
			},
			ProfessionalDocuments: physiotherapist.ProfessionalDocuments{
				License:     professionalLicense,
				Certificate: professionalCertificate,
			},
			PersonalDocuments: physiotherapist.PersonalDocuments{
				Nimc:    personalNimc,
				License: personalLicense,
			},
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
