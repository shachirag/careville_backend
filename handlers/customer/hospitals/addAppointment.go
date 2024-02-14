package hospitals

import (
	"time"

	"careville_backend/database"
	hospitals "careville_backend/dto/customer/hospitals"
	customerMiddleware "careville_backend/dto/customer/middleware"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Add appointment
// @Tags customer hospitals
// @Description Add appointment
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param serviceId query string true "service ID"
// @Param  customer body hospitals.HospitalClinicAppointmentReqDto true "add HospitalClinic"
// @Produce json
// @Success 200 {object} hospitals.HospitalClinicAppointmentResDto
// @Router /customer/healthFacility/add-hospClinic-appointment [post]
func AddHospClinicAppointment(c *fiber.Ctx) error {

	var (
		appointmentColl = database.GetCollection("appointment")
		serviceColl     = database.GetCollection("service")
		customerColl    = database.GetCollection("customer")
		data            hospitals.HospitalClinicAppointmentReqDto
		appointment     entity.AppointmentEntity
		service         entity.ServiceEntity
	)

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(hospitals.HospitalClinicAppointmentResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	familyObjectID, err := primitive.ObjectIDFromHex(data.FamillyMemberId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(hospitals.HospitalClinicAppointmentResDto{
			Status:  false,
			Message: "Invalid ID format",
		})
	}

	serviceId := c.Query("serviceId")

	if serviceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(hospitals.HospitalClinicAppointmentResDto{
			Status:  false,
			Message: "service Id is mandatory",
		})
	}

	serviceObjectID, err := primitive.ObjectIDFromHex(serviceId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(hospitals.HospitalClinicAppointmentResDto{
			Status:  false,
			Message: "Invalid ID format",
		})
	}

	doctorObjID, err := primitive.ObjectIDFromHex(data.DoctorId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(hospitals.HospitalClinicAppointmentResDto{
			Status:  false,
			Message: "Invalid ID format",
		})
	}

	doctorFilter := bson.M{
		"_id": serviceObjectID,
		"hospClinic.doctor": bson.M{
			"$elemMatch": bson.M{
				"id": doctorObjID,
			},
		},
	}

	doctorProjection := bson.M{
		"hospClinic.doctor.id":         1,
		"hospClinic.doctor.name":       1,
		"hospClinic.doctor.speciality": 1,
		"hospClinic.doctor.image":      1,
	}

	doctorOpts := options.FindOne().SetProjection(doctorProjection)

	err = serviceColl.FindOne(ctx, doctorFilter, doctorOpts).Decode(&service)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(hospitals.HospitalClinicAppointmentResDto{
			Status:  false,
			Message: "Failed to fetch doctor data: " + err.Error(),
		})
	}

	if service.HospClinic == nil {
		return c.Status(fiber.StatusNotFound).JSON(hospitals.HospitalClinicAppointmentResDto{
			Status:  false,
			Message: "Hospital clinic data not found",
		})
	}

	var doctorData entity.Doctor
	if service.HospClinic != nil && len(service.HospClinic.Doctor) > 0 {
		for _, doctor := range service.HospClinic.Doctor {
			if doctor.Id == doctorObjID {
				doctorData = doctor
				break
			}
		}
	}

	customerMiddlewareData := customerMiddleware.GetCustomerMiddlewareData(c)
	familyFilter := bson.M{
		"_id": customerMiddlewareData.CustomerId,
		"familyMembers": bson.M{
			"$elemMatch": bson.M{
				"id": familyObjectID,
			},
		},
	}

	var family entity.CustomerEntity
	familyProjection := bson.M{
		"familyMembers.id":           1,
		"familyMembers.name":         1,
		"familyMembers.age":          1,
		"familyMembers.sex":          1,
		"familyMembers.relationShip": 1,
	}

	familyOpts := options.FindOne().SetProjection(familyProjection)
	err = customerColl.FindOne(ctx, familyFilter, familyOpts).Decode(&family)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(hospitals.HospitalClinicAppointmentResDto{
			Status:  false,
			Message: "Failed to fetch family data: " + err.Error(),
		})
	}

	var familyData entity.FamilyMembers
	if family.FamilyMembers != nil {
		for _, family := range family.FamilyMembers {
			if family.Id == familyObjectID {
				familyData = family
				break
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
		return c.Status(fiber.StatusInternalServerError).JSON(hospitals.HospitalClinicAppointmentResDto{
			Status:  false,
			Message: "Failed to fetch customer data: " + err.Error(),
		})
	}

	var fromDate time.Time
	if data.FromDate != "" {
		fromDate, err = time.Parse(time.RFC3339, data.FromDate)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(hospitals.HospitalClinicAppointmentResDto{
				Status:  false,
				Message: "Failed to parse fromDate date: " + err.Error(),
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(hospitals.HospitalClinicAppointmentResDto{
			Status:  false,
			Message: "fromDate is mandatory",
		})
	}

	var toDate time.Time
	if data.ToDate != "" {
		toDate, err = time.Parse(time.RFC3339, data.ToDate)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(hospitals.HospitalClinicAppointmentResDto{
				Status:  false,
				Message: "Failed to parse toDate date: " + err.Error(),
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(hospitals.HospitalClinicAppointmentResDto{
			Status:  false,
			Message: "toDate date is mandatory",
		})
	}

	var remindMeBefore time.Time
	if data.RemindMeBefore != "" {
		remindMeBefore, err = time.Parse(time.RFC3339, data.RemindMeBefore)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(hospitals.HospitalClinicAppointmentResDto{
				Status:  false,
				Message: "Failed to parse remindMeBefore date: " + err.Error(),
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(hospitals.HospitalClinicAppointmentResDto{
			Status:  false,
			Message: "remindMeBefore date is mandatory",
		})
	}

	appointmentData := entity.HospitalAppointmentEntity{
		Doctor: entity.DoctorAppointmentEntity{
			ID:         doctorObjID,
			Name:       doctorData.Name,
			Image:      doctorData.Image,
			Speciality: doctorData.Speciality,
		},
		AppointmentDetails: entity.AppointmentDetailsAppointmentEntity{
			RemindMeBefore: remindMeBefore,
			From:           fromDate,
			To:             toDate,
		},
		FamilyMember: entity.FamilyMemberAppointmentEntity{
			ID:           familyObjectID,
			Name:         familyData.Name,
			Age:          familyData.Age,
			Sex:          familyData.Sex,
			Relationship: familyData.RelationShip,
		},
		FamilyType: data.FamilyType,
		PricePaid:  0,
	}

	appointment = entity.AppointmentEntity{
		Id:                   primitive.NewObjectID(),
		Role:                 "healthFacility",
		FacilityOrProfession: "hospClinic",
		ServiceID:            serviceObjectID,
		Customer: entity.CustomerAppointmentEntity{
			ID:          customerMiddlewareData.CustomerId,
			FirstName:   customer.FirstName,
			LastName:    customer.LastName,
			Image:       customer.Image,
			Email:       customer.Email,
			PhoneNumber: customer.PhoneNumber,
		},
		HospitalClinic:    &appointmentData,
		PaymentStatus:     "initiated",
		AppointmentStatus: "pending",
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	session, err := database.GetMongoClient().StartSession()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(hospitals.HospitalClinicAppointmentResDto{
			Status:  false,
			Message: "Failed to start session",
		})
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		_, err := appointmentColl.InsertOne(sessCtx, appointment)
		if err != nil {
			return nil, err
		}

		filter := bson.M{
			"_id": serviceObjectID,
			"hospClinic.doctor": bson.M{
				"$elemMatch": bson.M{
					"id": doctorObjID,
				},
			},
		}

		update := bson.M{
			"$push": bson.M{
				"hospClinic.doctor.$.upcommingEvents": bson.M{
					"id":        appointment.Id,
					"startTime": appointment.HospitalClinic.AppointmentDetails.From,
					"endTime":   appointment.HospitalClinic.AppointmentDetails.To,
				},
			},
		}

		_, err = serviceColl.UpdateOne(sessCtx, filter, update)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	_, err = session.WithTransaction(ctx, callback)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(hospitals.HospitalClinicAppointmentResDto{
			Status:  false,
			Message: "Failed to update appointment data: " + err.Error(),
		})
	}

	hospClinicRes := hospitals.HospitalClinicAppointmentResDto{
		Status:  true,
		Message: "Appointment added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(hospClinicRes)
}
