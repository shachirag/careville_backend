package hospitals

import (
	"careville_backend/database"
	services "careville_backend/dto/admin/services/hospitals"
	"careville_backend/entity"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary get-hospital detail 
// @Description get-hospital detail 
// @Tags admin hospitals
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} services.GetHospitalDetailResDto
// @Router /admin/healthFacility/get-hospital/{id} [get]
func GetHospitalDetail(c *fiber.Ctx) error {
	ctx := context.TODO()

	var service entity.ServiceEntity

	serviceColl := database.GetCollection("service")

	idParam := c.Params("id")
	hospitalID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.GetHospitalDetailResDto{
			Status:  false,
			Message: "Invalid appointment ID",
		})
	}

	filter := bson.M{
		"_id": hospitalID,
	}

	projection := bson.M{
		"_id":                                   1,
		"facilityOrProfession":                  1,
		"role":                                  1,
		"profileId":                             1,
		"serviceStatus":                         1,
		"user.firstName":                        1,
		"user.lastName":                         1,
		"user.email":                            1,
		"user.phoneNumber.dialCode":             1,
		"user.phoneNumber.countryCode":          1,
		"user.phoneNumber.number":               1,
		"hospClinic.documents.certificate":      1,
		"hospClinic.documents.license":          1,
		"hospClinic.information.name":           1,
		"hospClinic.information.image":          1,
		"hospClinic.information.additionalText": 1,
		"hospClinic.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
		"hospClinic.doctor.id":         1,
		"hospClinic.doctor.name":       1,
		"hospClinic.doctor.speciality": 1,
		"hospClinic.doctor.schedule": bson.M{
			"startTime": 1,
			"endTime":   1,
			"days":      1,
		},
		"hospClinic.otherServices": 1,
		"hospClinic.insurances":    1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.GetHospitalDetailResDto{
				Status:  false,
				Message: "Other service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.GetHospitalDetailResDto{
			Status:  false,
			Message: "Failed to fetch other service from MongoDB: " + err.Error(),
		})
	}

	if service.HospClinic == nil {
		return c.Status(fiber.StatusOK).JSON(services.GetHospitalDetailResDto{
			Status:  false,
			Message: "HospClinic information not found.",
		})
	}

	var doctorData []services.Doctor
	if service.HospClinic != nil && len(service.HospClinic.Doctor) > 0 {
		for _, doctor := range service.HospClinic.Doctor {
			scheduleData := make([]services.Schedule, 0)
			for _, schedule := range doctor.Schedule {
				scheduleData = append(scheduleData, services.Schedule{
					StartTime: schedule.StartTime,
					EndTime:   schedule.EndTime,
					Days:      schedule.Days,
				})
			}
			doctorData = append(doctorData, services.Doctor{
				Id:         doctor.Id,
				Name:       doctor.Name,
				Speciality: doctor.Speciality,
				Schedule:   scheduleData,
			})
		}
	}

	var hospitalServices []string
	var hospitalInsurances []string
	var hospitalImage string
	var hospitalName string
	var additionalText string
	var hospitalAddress services.Address
	if service.HospClinic != nil {
		hospitalServices = service.HospClinic.OtherServices
		hospitalInsurances = service.HospClinic.Insurances
		hospitalName = service.HospClinic.Information.Name
		hospitalImage = service.HospClinic.Information.Image
		additionalText = service.HospClinic.Information.AdditionalText
		hospitalAddress = services.Address(service.HospClinic.Information.Address)
	}

	response := services.GetHospitalDetailResDto{
		Status:  true,
		Message: "Data fetched retrieved successfully",
		Data: services.GetHospitalDetailRes{
			Id:                   service.Id,
			FacilityOrProfession: service.FacilityOrProfession,
			Role:                 service.Role,
			ProfileId:            service.ProfileId,
			ServiceStatus:        service.ServiceStatus,
			User: services.User{
				FirstName: service.User.FirstName,
				LastName:  service.User.LastName,
				Email:     service.User.Email,
				PhoneNumber: services.PhoneNumber{
					DialCode: service.User.PhoneNumber.DialCode,
					Number:   service.User.PhoneNumber.Number,
				},
			},
			Doctor: doctorData,
			HospitalInformation: services.HospitalInformation{
				Name:           hospitalName,
				Image:          hospitalImage,
				AdditionalText: additionalText,
				Address:        hospitalAddress,
			},
			OtherServices: hospitalServices,
			Insurances:    hospitalInsurances,
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
