package emergency

import (
	"time"

	"careville_backend/database"
	"careville_backend/dto/customer/doctorProfession"
	customerMiddleware "careville_backend/dto/customer/middleware"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Add emergency doctor
// @Tags customer emergency
// @Description Add emergency doctor
// @Accept multipart/form-data
//
// @Param Authorization header string true "Authentication header"
//
// @Param serviceId query string true "service ID"
// @Param  customer body doctorProfession.DoctorProfessionAppointmentReqDto true "add doctorProfession"
// @Produce json
// @Success 200 {object} doctorProfession.AddEmergencyDoctorResDto
// @Router /customer/healthProfessional/emergency/add-emergency-doctor [post]
func AddEmergencyDoctor(c *fiber.Ctx) error {

	var (
		customerColl  = database.GetCollection("customer")
		emergencyColl = database.GetCollection("emergency")
		serviceColl   = database.GetCollection("service")
		data          doctorProfession.AddEmergencyDoctorReqDto
		emergency     entity.EmergencyEntity
	)

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(doctorProfession.AddEmergencyDoctorResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	serviceId := c.Query("serviceId")

	if serviceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(doctorProfession.AddEmergencyDoctorResDto{
			Status:  false,
			Message: "service Id is mandatory",
		})
	}

	serviceObjectID, err := primitive.ObjectIDFromHex(serviceId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(doctorProfession.AddEmergencyDoctorResDto{
			Status:  false,
			Message: "Invalid ID format",
		})
	}

	customerMiddlewareData := customerMiddleware.GetCustomerMiddlewareData(c)

	var customer entity.CustomerEntity
	customerFilter := bson.M{
		"_id": customerMiddlewareData.CustomerId,
	}

	customerProjection := bson.M{
		"_id":       1,
		"firstName": 1,
		"lastName":  1,
		"image":     1,
		"email":     1,
		"phoneNumber": bson.M{
			"dialCode":    1,
			"number":      1,
			"countryCode": 1,
		},
		"address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
	}

	customerOpts := options.FindOne().SetProjection(customerProjection)
	err = customerColl.FindOne(ctx, customerFilter, customerOpts).Decode(&customer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(doctorProfession.AddEmergencyDoctorResDto{
			Status:  false,
			Message: "Failed to fetch customer data: " + err.Error(),
		})
	}

	serviceFilter := bson.M{
		"_id":                  serviceObjectID,
		"facilityOrProfession": "doctor",
		"role":                 "healthProfessional",
	}

	serviceProjection := bson.M{
		"doctor.information.name":            1,
		"doctor.information.image":           1,
		"doctor.addionalServices.speciality": 1,
	}

	serviceOpts := options.FindOne().SetProjection(serviceProjection)

	var service entity.ServiceEntity
	err = serviceColl.FindOne(ctx, serviceFilter, serviceOpts).Decode(&service)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(doctorProfession.AddEmergencyDoctorResDto{
			Status:  false,
			Message: "Failed to fetch Doctor data: " + err.Error(),
		})
	}

	if service.Doctor == nil {
		return c.Status(fiber.StatusNotFound).JSON(doctorProfession.AddEmergencyDoctorResDto{
			Status:  false,
			Message: "Doctor data not found",
		})
	}

	var name string
	var image string
	var speciality string

	if service.Doctor != nil {
		name = service.Doctor.Information.Name
		image = service.Doctor.Information.Image
		speciality = service.Doctor.AdditionalServices.Speciality
	}

	emergency = entity.EmergencyEntity{
		ID:                   primitive.NewObjectID(),
		Role:                 "healthProfessional",
		FacilityOrProfession: "doctor",
		ServiceID:            serviceObjectID,
		Customer: entity.CustomerEmergencyEntity{
			ID:          customerMiddlewareData.CustomerId,
			FirstName:   customer.FirstName,
			LastName:    customer.LastName,
			Image:       customer.Image,
			Email:       customer.Email,
			PhoneNumber: customer.PhoneNumber,
			Address: entity.Address{
				Coordinates: customer.Address.Coordinates,
				Type:        customer.Address.Type,
				Add:         customer.Address.Add,
			},
		},
		Doctor: entity.DoctorEmergencyEntity{
			ID:         serviceObjectID,
			Name:       name,
			Image:      image,
			Speciality: speciality,
		},
		Price: entity.PriceEmergencyEntity{
			PricePaid: data.PricePaid,
		},
		Type:      "EmergencyDoctor",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	_, err = emergencyColl.InsertOne(ctx, emergency)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(doctorProfession.AddEmergencyDoctorResDto{
			Status:  false,
			Message: "Failed to insert emergency doctor data into MongoDB: " + err.Error(),
		})
	}

	doctorRes := doctorProfession.AddEmergencyDoctorResDto{
		Status:  true,
		Message: "Successfully added",
	}
	return c.Status(fiber.StatusOK).JSON(doctorRes)
}
