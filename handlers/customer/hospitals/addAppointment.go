package hospitals

import (
	"encoding/json"
	"time"

	"careville_backend/database"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"
	"careville_backend/entity/subEntity"

	"github.com/gofiber/fiber/v2"
)

// @Summary Add appointment
// @Tags customer hospitals
// @Description Add appointment
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param  provider body services.HospitalClinicRequestDto true "add HospitalClinic"
// @Produce json
// @Success 200 {object} services.HospitalClinicResDto
// @Router /customer/healthFacility/add-appointment [post]
func AddHospClinicAppointment(c *fiber.Ctx) error {
	var (
		appointmentColl = database.GetCollection("appointment")
		data            services.HospitalClinicRequestDto
		appointment     entity.AppointmentEntity
	)

	dataStr := c.FormValue("data")
	dataBytes := []byte(dataStr)

	err := json.Unmarshal(dataBytes, &data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.HospitalClinicResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	hospClinicData := subEntity.HospClinicUpdateServiceSubEntity{
		Information: subEntity.InformationUpdateServiceSubEntity{
			Name:           data.HospitalClinicReqDto.InformationName,
			AdditionalText: data.HospitalClinicReqDto.AdditionalText,
			Image:          hospitalImageUrl,
			Address: subEntity.AddressUpdateServiceSubEntity{
				Coordinates: []float64{longitude, latitude},
				Add:         data.HospitalClinicReqDto.Address,
				Type:        "Point",
			},
			IsEmergencyAvailable: false,
		},
		Documents: subEntity.DocumentsUpdateServiceSubEntity{
			Certificate: hospitalCertificateUrl,
			License:     hospitalLicenceUrl,
		},
		OtherServices: data.HospitalClinicReqDto.OtherServices,
		Insurances:    data.HospitalClinicReqDto.Insurances,
		Doctor:        doctors,
	}

	appointment = entity.AppointmentEntity{
		Role:                 "healthFacility",
		FacilityOrProfession: "hospClinic",
		ServiceStatus:        "pending",
		HospClinic:           &hospClinicData,
		UpdatedAt:            time.Now().UTC(),
	}

	_, err = appointmentColl.InsertOne(ctx, appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.HospitalClinicResDto{
			Status:  false,
			Message: "Failed to insert hospital/clinic data into MongoDB: " + err.Error(),
		})
	}

	hospClinicRes := services.HospitalClinicResDto{
		Status:  true,
		Message: "Appointment added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(hospClinicRes)
}
