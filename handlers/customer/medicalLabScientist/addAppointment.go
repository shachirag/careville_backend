package medicalLabScientist

import (
	"time"

	"careville_backend/database"
	"careville_backend/dto/customer/hospitals"
	"careville_backend/dto/customer/medicalLabScientist"
	customerMiddleware "careville_backend/dto/customer/middleware"
	"careville_backend/entity"
	"careville_backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Add appointment
// @Tags customer medicalLabScientist
// @Description Add appointment
// @Accept multipart/form-data
//
// @Param Authorization header string true "Authentication header"
//
// @Param serviceId query string true "service ID"
// @Param  customer body  medicalLabScientist.MedicalLabScientistAppointmentReqDto true "add medicalLabScientist"
// @Produce json
// @Success 200 {object}  medicalLabScientist.MedicalLabScientistAppointmentResDto
// @Router /customer/healthProfessional/add-medicalLabScientist-appointment [post]
func AddMedicalLabScientistAppointment(c *fiber.Ctx) error {

	var (
		appointmentColl = database.GetCollection("appointment")
		customerColl    = database.GetCollection("customer")
		serviceColl     = database.GetCollection("service")
		data            medicalLabScientist.MedicalLabScientistAppointmentReqDto
		appointment     entity.AppointmentEntity
	)

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(medicalLabScientist.MedicalLabScientistAppointmentResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	var familyObjectID primitive.ObjectID

	if data.FamillyMemberId != nil && *data.FamillyMemberId != "" {
		familyObjectID, err = primitive.ObjectIDFromHex(*data.FamillyMemberId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(medicalLabScientist.MedicalLabScientistAppointmentResDto{
				Status:  false,
				Message: "Invalid ID format",
			})
		}
	}
	serviceId := c.Query("serviceId")

	if serviceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(medicalLabScientist.MedicalLabScientistAppointmentResDto{
			Status:  false,
			Message: "service Id is mandatory",
		})
	}

	serviceObjectID, err := primitive.ObjectIDFromHex(serviceId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(medicalLabScientist.MedicalLabScientistAppointmentResDto{
			Status:  false,
			Message: "Invalid ID format",
		})
	}

	medicalLabScientistServiceDataServiceObjID, err := primitive.ObjectIDFromHex(data.MedicalLabScientistServiceId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(medicalLabScientist.MedicalLabScientistAppointmentResDto{
			Status:  false,
			Message: "Invalid ID format",
		})
	}

	var fromDate time.Time
	if data.FromDate != "" {
		fromDate, err = time.Parse(time.DateTime, data.FromDate)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.MedicalLabScientistAppointmentResDto{
				Status:  false,
				Message: "Failed to parse fromDate date: " + err.Error(),
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(medicalLabScientist.MedicalLabScientistAppointmentResDto{
			Status:  false,
			Message: "fromDate is mandatory",
		})
	}

	var toDate time.Time
	if data.ToDate != "" {
		toDate, err = time.Parse(time.DateTime, data.ToDate)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.MedicalLabScientistAppointmentResDto{
				Status:  false,
				Message: "Failed to parse toDate date: " + err.Error(),
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(medicalLabScientist.MedicalLabScientistAppointmentResDto{
			Status:  false,
			Message: "toDate date is mandatory",
		})
	}

	customerMiddlewareData := customerMiddleware.GetCustomerMiddlewareData(c)

	overlapFilter := bson.M{
		"customer.id":                    customerMiddlewareData.CustomerId,
		"serviceId":                      serviceObjectID,
		"medicalLabScientist.service.id": medicalLabScientistServiceDataServiceObjID,
		"medicalLabScientist.appointmentDetails.from": bson.M{"$lte": toDate},
		"medicalLabScientist.appointmentDetails.to":   bson.M{"$gte": fromDate},
	}

	count, err := appointmentColl.CountDocuments(ctx, overlapFilter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.MedicalLabScientistAppointmentResDto{
			Status:  false,
			Message: "Failed to check existing appointments: " + err.Error(),
		})
	}

	if count > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(medicalLabScientist.MedicalLabScientistAppointmentResDto{
			Status:  false,
			Message: "You have already created a booking for this service with the Medical Lab Scientist.",
		})
	}

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
			return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.MedicalLabScientistAppointmentResDto{
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
		"age": 1,
	}

	customerOpts := options.FindOne().SetProjection(customerProjection)
	err = customerColl.FindOne(ctx, customerFilter, customerOpts).Decode(&customer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.MedicalLabScientistAppointmentResDto{
			Status:  false,
			Message: "Failed to fetch customer data: " + err.Error(),
		})
	}

	serviceFilter := bson.M{
		"_id":                  serviceObjectID,
		"facilityOrProfession": "medicalLabScientist",
		"role":                 "healthProfessional",
	}

	serviceProjection := bson.M{
		"_id":                                                1,
		"user.notification.deviceToken":                      1,
		"user.notification.deviceType":                       1,
		"medicalLabScientist.information.name":               1,
		"medicalLabScientist.information.image":              1,
		"medicalLabScientist.professionalDetails.department": 1,
	}

	serviceOpts := options.FindOne().SetProjection(serviceProjection)

	var service entity.ServiceEntity
	err = serviceColl.FindOne(ctx, serviceFilter, serviceOpts).Decode(&service)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.MedicalLabScientistAppointmentResDto{
			Status:  false,
			Message: "Failed to fetch MedicalLabScientist data: " + err.Error(),
		})
	}

	if service.MedicalLabScientist == nil {
		return c.Status(fiber.StatusNotFound).JSON(medicalLabScientist.MedicalLabScientistAppointmentResDto{
			Status:  false,
			Message: "MedicalLabScientist data not found",
		})
	}

	medicalLabScientistServiceFilter := bson.M{
		"_id": serviceObjectID,
		"medicalLabScientist.serviceAndSchedule": bson.M{
			"$elemMatch": bson.M{
				"id": medicalLabScientistServiceDataServiceObjID,
			},
		},
	}

	medicalLabScientistProjection := bson.M{
		"medicalLabScientist.serviceAndSchedule.id":          1,
		"medicalLabScientist.serviceAndSchedule.name":        1,
		"medicalLabScientist.serviceAndSchedule.serviceFees": 1,
	}

	medicalLabScientistServiceOpts := options.FindOne().SetProjection(medicalLabScientistProjection)

	err = serviceColl.FindOne(ctx, medicalLabScientistServiceFilter, medicalLabScientistServiceOpts).Decode(&service)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.MedicalLabScientistAppointmentResDto{
			Status:  false,
			Message: "Failed to fetch service data: " + err.Error(),
		})
	}

	var medicalLabScientistServiceData entity.ServiceAndSchedule
	if service.MedicalLabScientist != nil && len(service.MedicalLabScientist.ServiceAndSchedule) > 0 {
		for _, medicalLabScientistService := range service.MedicalLabScientist.ServiceAndSchedule {
			if medicalLabScientistService.Id == medicalLabScientistServiceDataServiceObjID {
				medicalLabScientistServiceData = medicalLabScientistService
				break
			}
		}
	}

	var remindMeBefore time.Time
	if data.RemindMeBefore != "" {
		remindMeBefore, err = time.Parse(time.DateTime, data.RemindMeBefore)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.MedicalLabScientistAppointmentResDto{
				Status:  false,
				Message: "Failed to parse remindMeBefore date: " + err.Error(),
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(medicalLabScientist.MedicalLabScientistAppointmentResDto{
			Status:  false,
			Message: "remindMeBefore date is mandatory",
		})
	}

	var name string
	var image string
	var department string

	if service.MedicalLabScientist != nil {
		name = service.MedicalLabScientist.Information.Name
		image = service.MedicalLabScientist.Information.Image
		department = service.MedicalLabScientist.ProfessionalDetails.Department
	}

	appointmentData := entity.MedicalLabScientistAppointmentEntity{
		AppointmentDetails: entity.AppointmentDetailsAppointmentEntity{
			From:           fromDate,
			To:             toDate,
			RemindMeBefore: remindMeBefore,
		},
		Service: entity.ServiceAppointmentEntity{
			Id:          medicalLabScientistServiceData.Id,
			Name:        medicalLabScientistServiceData.Name,
			ServiceFees: medicalLabScientistServiceData.ServiceFees,
		},
		Information: entity.MedicalLabScientistInformation{
			Name:       name,
			Image:      image,
			Department: department,
		},
		FamilyMember: entity.FamilyMemberAppointmentEntity{
			ID:           familyObjectID,
			Name:         familyData.Name,
			Age:          familyData.Age,
			Sex:          familyData.Sex,
			Relationship: familyData.RelationShip,
		},
		FamilyType: data.FamilyType,
		PricePaid:  data.PricePaid,
	}

	var id = primitive.NewObjectID()
	appointment = entity.AppointmentEntity{
		Id:                   id,
		Role:                 "healthProfessional",
		FacilityOrProfession: "medicalLabScientist",
		ServiceID:            serviceObjectID,
		Customer: entity.CustomerAppointmentEntity{
			ID:          customerMiddlewareData.CustomerId,
			FirstName:   customer.FirstName,
			LastName:    customer.LastName,
			Image:       customer.Image,
			Email:       customer.Email,
			PhoneNumber: customer.PhoneNumber,
			Age:         customer.Age,
		},
		MedicalLabScientist: &appointmentData,
		PaymentStatus:       "initiated",
		AppointmentStatus:   "pending",
		CreatedAt:           time.Now().UTC(),
		UpdatedAt:           time.Now().UTC(),
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
			"_id":                  serviceObjectID,
			"facilityOrProfession": "medicalLabScientist",
		}

		update := bson.M{
			"$push": bson.M{
				"medicalLabScientist.upcommingEvents": bson.M{
					"id":        id,
					"startTime": fromDate,
					"endTime":   toDate,
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

	if service.User.Notification.DeviceToken != "" && service.User.Notification.DeviceType != "" {
		formattedFromDate := fromDate.Format("02 Jan 2006 15:04")
		formattedToDate := toDate.Format("02 Jan 2006 15:04")

		notificationTitle := "Appointment Scheduled for Medical Lab Scientist"
		notificationBody := "An appointment has been scheduled for a medical lab scientist from " + formattedFromDate + " to " + formattedToDate + "."
		notificationData := map[string]string{
			"type":                 "appointment-scheduled",
			"appointmentId":        id.Hex(),
			"role":                 "healthProfessional",
			"facilityOrProfession": "medicalLabScientist",
		}

		utils.SendNotificationToUser(service.User.Notification.DeviceToken, service.User.Notification.DeviceType, notificationTitle, notificationBody, notificationData, service.Id, "provider")
	}

	medicalLabScientistRes := medicalLabScientist.MedicalLabScientistAppointmentResDto{
		Status:  true,
		Message: "Appointment added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(medicalLabScientistRes)
}
