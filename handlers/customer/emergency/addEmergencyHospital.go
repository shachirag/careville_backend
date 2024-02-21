package emergency

import (
	"strconv"
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
// @Param  customer body doctorProfession.DoctorProfessionAppointmentReqDto true "add doctorProfession"
// @Produce json
// @Success 200 {object} doctorProfession.AddEmergencyDoctorResDto
// @Router /customer/emergency/add-emergency-hospital [post]
func AddEmergencyHospital(c *fiber.Ctx) error {

	var (
		customerColl  = database.GetCollection("customer")
		emergencyColl = database.GetCollection("emergency")
		serviceColl   = database.GetCollection("service")
		data          doctorProfession.AddEmergencyHospitalReqDto
		emergency     entity.EmergencyEntity
	)

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(doctorProfession.AddEmergencyDoctorResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	if data.HospitalId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(doctorProfession.AddEmergencyDoctorResDto{
			Status:  false,
			Message: "hospital Id is mandatory",
		})
	}

	hospitalObjectID, err := primitive.ObjectIDFromHex(data.HospitalId)
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
		"_id":                  hospitalObjectID,
		"facilityOrProfession": "hospClinic",
		"role":                 "healthFacility",
	}

	serviceProjection := bson.M{
		"hospClinic.information.name":  1,
		"hospClinic.information.image": 1,
		"hospClinic.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
	}

	serviceOpts := options.FindOne().SetProjection(serviceProjection)

	var service entity.ServiceEntity
	err = serviceColl.FindOne(ctx, serviceFilter, serviceOpts).Decode(&service)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(doctorProfession.AddEmergencyDoctorResDto{
			Status:  false,
			Message: "Failed to fetch Hospital data: " + err.Error(),
		})
	}

	if service.HospClinic == nil {
		return c.Status(fiber.StatusNotFound).JSON(doctorProfession.AddEmergencyDoctorResDto{
			Status:  false,
			Message: "Hospital data not found",
		})
	}

	longitude, err := strconv.ParseFloat(data.Longitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(doctorProfession.AddEmergencyDoctorResDto{
			Status:  false,
			Message: "Invalid longitude format",
		})
	}

	latitude, err := strconv.ParseFloat(data.Latitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(doctorProfession.AddEmergencyDoctorResDto{
			Status:  false,
			Message: "Invalid latitude format",
		})
	}

	var name string
	var image string
	var address entity.Address

	if service.HospClinic != nil {
		name = service.HospClinic.Information.Name
		image = service.HospClinic.Information.Image
		address = service.HospClinic.Information.Address
	}

	hospitalData := entity.HospitalEmergencyEntity{
		Information: entity.HospitalInformation{
			ID:      hospitalObjectID,
			Name:    name,
			Image:   image,
			Address: address,
		},
		AddedAddress: entity.Address{
			Coordinates: []float64{longitude, latitude},
			Add:         data.Address,
			Type:        "Point",
		},
	}

	emergency = entity.EmergencyEntity{
		ID:                   primitive.NewObjectID(),
		Role:                 "healthProfessional",
		FacilityOrProfession: "doctor",
		ServiceID:            hospitalObjectID,
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

		Hospital: &hospitalData,
		Price: entity.PriceEmergencyEntity{
			PricePaid: 0,
		},
		Type:      "EmergencyHospital",
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
