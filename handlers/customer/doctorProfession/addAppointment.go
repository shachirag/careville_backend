package doctorProfession

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

// @Summary Add appointment
// @Tags customer doctorProfession
// @Description Add appointment
// @Accept multipart/form-data
//
// @Param Authorization header string true "Authentication header"
//
// @Param serviceId query string true "service ID"
// @Param  customer body doctorProfession.DoctorProfessionAppointmentReqDto true "add doctorProfession"
// @Produce json
// @Success 200 {object} doctorProfession.DoctorProfessionAppointmentResDto
// @Router /customer/healthProfessional/add-doctor-appointment [post]
func AddDoctorAppointment(c *fiber.Ctx) error {

	var (
		appointmentColl = database.GetCollection("appointment")
		customerColl    = database.GetCollection("customer")
		serviceColl     = database.GetCollection("service")
		data            doctorProfession.DoctorProfessionAppointmentReqDto
		appointment     entity.AppointmentEntity
	)

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(doctorProfession.DoctorProfessionAppointmentResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	var familyObjectID primitive.ObjectID

	if data.FamillyMemberId != nil && *data.FamillyMemberId != "" {

		familyObjectID, err = primitive.ObjectIDFromHex(*data.FamillyMemberId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(doctorProfession.DoctorProfessionAppointmentResDto{
				Status:  false,
				Message: "Invalid ID format",
			})
		}
	}

	serviceId := c.Query("serviceId")

	if serviceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(doctorProfession.DoctorProfessionAppointmentResDto{
			Status:  false,
			Message: "service Id is mandatory",
		})
	}

	serviceObjectID, err := primitive.ObjectIDFromHex(serviceId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(doctorProfession.DoctorProfessionAppointmentResDto{
			Status:  false,
			Message: "Invalid ID format",
		})
	}

	customerMiddlewareData := customerMiddleware.GetCustomerMiddlewareData(c)
	var familyData entity.FamilyMembers

	if data.FamillyMemberId != nil && *data.FamillyMemberId != "" {

		familyFilter := bson.M{
			"_id": customerMiddlewareData.CustomerId,
			"familyMembers": bson.M{
				"$elemMatch": bson.M{
					"id": familyObjectID,
				},
			},
		}

		familyProjection := bson.M{
			"familyMembers.id":           1,
			"familyMembers.name":         1,
			"familyMembers.age":          1,
			"familyMembers.sex":          1,
			"familyMembers.relationShip": 1,
		}

		familyOpts := options.FindOne().SetProjection(familyProjection)

		var family entity.CustomerEntity
		err = customerColl.FindOne(ctx, familyFilter, familyOpts).Decode(&family)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(doctorProfession.DoctorProfessionAppointmentResDto{
				Status:  false,
				Message: "Failed to fetch family data: " + err.Error(),
			})
		}

		if family.FamilyMembers != nil {
			for _, family := range family.FamilyMembers {
				if family.Id == familyObjectID {
					familyData = family
					break
				}
			}
		}
	}

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
	}

	customerOpts := options.FindOne().SetProjection(customerProjection)
	err = customerColl.FindOne(ctx, customerFilter, customerOpts).Decode(&customer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(doctorProfession.DoctorProfessionAppointmentResDto{
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
		return c.Status(fiber.StatusInternalServerError).JSON(doctorProfession.DoctorProfessionAppointmentResDto{
			Status:  false,
			Message: "Failed to fetch Doctor data: " + err.Error(),
		})
	}

	if service.Doctor == nil {
		return c.Status(fiber.StatusNotFound).JSON(doctorProfession.DoctorProfessionAppointmentResDto{
			Status:  false,
			Message: "Doctor data not found",
		})
	}

	longitude, err := strconv.ParseFloat(data.Longitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(doctorProfession.DoctorProfessionAppointmentResDto{
			Status:  false,
			Message: "Invalid longitude format",
		})
	}

	latitude, err := strconv.ParseFloat(data.Latitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(doctorProfession.DoctorProfessionAppointmentResDto{
			Status:  false,
			Message: "Invalid latitude format",
		})
	}

	var fromDate time.Time
	if data.FromDate != "" {
		fromDate, err = time.Parse(time.DateTime, data.FromDate)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(doctorProfession.DoctorProfessionAppointmentResDto{
				Status:  false,
				Message: "Failed to parse fromDate date: " + err.Error(),
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(doctorProfession.DoctorProfessionAppointmentResDto{
			Status:  false,
			Message: "fromDate is mandatory",
		})
	}

	var toDate time.Time
	if data.ToDate != "" {
		toDate, err = time.Parse(time.DateTime, data.ToDate)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(doctorProfession.DoctorProfessionAppointmentResDto{
				Status:  false,
				Message: "Failed to parse toDate date: " + err.Error(),
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(doctorProfession.DoctorProfessionAppointmentResDto{
			Status:  false,
			Message: "toDate date is mandatory",
		})
	}

	var remindMeBefore time.Time
	if data.RemindMeBefore != "" {
		remindMeBefore, err = time.Parse(time.DateTime, data.RemindMeBefore)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(doctorProfession.DoctorProfessionAppointmentResDto{
				Status:  false,
				Message: "Failed to parse remindMeBefore date: " + err.Error(),
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(doctorProfession.DoctorProfessionAppointmentResDto{
			Status:  false,
			Message: "remindMeBefore date is mandatory",
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

	appointmentData := entity.DoctorProfessionAppointmentEntity{
		AppointmentDetails: entity.AppointmentDetailsAppointmentEntity{
			From:           fromDate,
			To:             toDate,
			RemindMeBefore: remindMeBefore,
		},
		Information: entity.DoctorProfessionInformation{
			Name:       name,
			Image:      image,
			Speciality: speciality,
		},
		FamilyMember: entity.FamilyMemberAppointmentEntity{
			ID:           familyObjectID,
			Name:         familyData.Name,
			Age:          familyData.Age,
			Sex:          familyData.Sex,
			Relationship: familyData.RelationShip,
		},
		Destination: entity.Address{
			Coordinates: []float64{longitude, latitude},
			Add:         data.Address,
			Type:        "Point",
		},
		FamilyType: data.FamilyType,
		PricePaid:  data.PricePaid,
	}

	appointment = entity.AppointmentEntity{
		Id:                   primitive.NewObjectID(),
		Role:                 "healthProfessional",
		FacilityOrProfession: "doctor",
		ServiceID:            serviceObjectID,
		Customer: entity.CustomerAppointmentEntity{
			ID:          customerMiddlewareData.CustomerId,
			FirstName:   customer.FirstName,
			LastName:    customer.LastName,
			Image:       customer.Image,
			Email:       customer.Email,
			PhoneNumber: customer.PhoneNumber,
		},
		Doctor:            &appointmentData,
		PaymentStatus:     "initiated",
		AppointmentStatus: "pending",
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	_, err = appointmentColl.InsertOne(ctx, appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(doctorProfession.DoctorProfessionAppointmentResDto{
			Status:  false,
			Message: "Failed to insert doctor appointment data into MongoDB: " + err.Error(),
		})
	}

	doctorRes := doctorProfession.DoctorProfessionAppointmentResDto{
		Status:  true,
		Message: "Appointment added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(doctorRes)
}
