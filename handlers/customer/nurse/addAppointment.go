package nurse

import (
	"strconv"
	"time"

	"careville_backend/database"
	customerMiddleware "careville_backend/dto/customer/middleware"
	"careville_backend/dto/customer/nurse"
	"careville_backend/entity"
	"careville_backend/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Add appointment
// @Tags customer nurse
// @Description Add appointment
// @Accept multipart/form-data
//
// @Param Authorization header string true "Authentication header"
//
// @Param serviceId query string true "service ID"
// @Param  customer body  nurse.NurseAppointmentReqDto true "add nurse"
// @Produce json
// @Success 200 {object}  nurse.AppoiynmentResDto
// @Router /customer/healthProfessional/add-nurse-appointment [post]
func AddNurseAppointment(c *fiber.Ctx) error {

	var (
		appointmentColl = database.GetCollection("appointment")
		customerColl    = database.GetCollection("customer")
		serviceColl     = database.GetCollection("service")
		data            nurse.NurseAppointmentReqDto
		appointment     entity.AppointmentEntity
	)

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(nurse.AppoiynmentResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	customerMiddlewareData := customerMiddleware.GetCustomerMiddlewareData(c)

	var familyObjectID primitive.ObjectID

	if data.FamillyMemberId != nil && *data.FamillyMemberId != "" {

		familyObjectID, err = primitive.ObjectIDFromHex(*data.FamillyMemberId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(nurse.AppoiynmentResDto{
				Status:  false,
				Message: "Invalid ID format",
			})
		}
	}

	serviceId := c.Query("serviceId")

	if serviceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(nurse.AppoiynmentResDto{
			Status:  false,
			Message: "service Id is mandatory",
		})
	}

	serviceObjectID, err := primitive.ObjectIDFromHex(serviceId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(nurse.AppoiynmentResDto{
			Status:  false,
			Message: "Invalid ID format",
		})
	}

	nurseServiceDataServiceObjID, err := primitive.ObjectIDFromHex(data.NurseServiceId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(nurse.AppoiynmentResDto{
			Status:  false,
			Message: "Invalid ID format",
		})
	}

	var fromDate time.Time
	if data.FromDate != "" {
		fromDate, err = time.Parse(time.DateTime, data.FromDate)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(nurse.AppoiynmentResDto{
				Status:  false,
				Message: "Failed to parse fromDate date: " + err.Error(),
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(nurse.AppoiynmentResDto{
			Status:  false,
			Message: "fromDate is mandatory",
		})
	}

	var toDate time.Time
	if data.ToDate != "" {
		toDate, err = time.Parse(time.DateTime, data.ToDate)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(nurse.AppoiynmentResDto{
				Status:  false,
				Message: "Failed to parse toDate date: " + err.Error(),
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(nurse.AppoiynmentResDto{
			Status:  false,
			Message: "toDate date is mandatory",
		})
	}

	overlapFilter := bson.M{
		"customer.id":                   customerMiddlewareData.CustomerId,
		"serviceId":                     serviceObjectID,
		"nurse.service.id":              nurseServiceDataServiceObjID,
		"nurse.appointmentDetails.from": bson.M{"$lte": toDate},
		"nurse.appointmentDetails.to":   bson.M{"$gte": fromDate},
	}

	count, err := appointmentColl.CountDocuments(ctx, overlapFilter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(nurse.AppoiynmentResDto{
			Status:  false,
			Message: "Failed to check existing appointments: " + err.Error(),
		})
	}

	if count > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(nurse.AppoiynmentResDto{
			Status:  false,
			Message: "You have already created a booking for this service with the nurse.",
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
			return c.Status(fiber.StatusInternalServerError).JSON(nurse.AppoiynmentResDto{
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
		return c.Status(fiber.StatusInternalServerError).JSON(nurse.AppoiynmentResDto{
			Status:  false,
			Message: "Failed to fetch customer data: " + err.Error(),
		})
	}

	serviceFilter := bson.M{
		"_id":                  serviceObjectID,
		"facilityOrProfession": "nurse",
		"role":                 "healthProfessional",
	}

	serviceProjection := bson.M{
		"_id":                           1,
		"user.notification.deviceToken": 1,
		"user.notification.deviceType":  1,
		"nurse.information.name":        1,
		"nurse.information.image":       1,
	}

	serviceOpts := options.FindOne().SetProjection(serviceProjection)

	var service entity.ServiceEntity
	err = serviceColl.FindOne(ctx, serviceFilter, serviceOpts).Decode(&service)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(nurse.AppoiynmentResDto{
			Status:  false,
			Message: "Failed to fetch nurse data: " + err.Error(),
		})
	}

	if service.Nurse == nil {
		return c.Status(fiber.StatusNotFound).JSON(nurse.AppoiynmentResDto{
			Status:  false,
			Message: "nurse data not found",
		})
	}

	nurseServiceFilter := bson.M{
		"_id": serviceObjectID,
		"nurse.schedule": bson.M{
			"$elemMatch": bson.M{
				"id": nurseServiceDataServiceObjID,
			},
		},
	}

	nurseProjection := bson.M{
		"nurse.schedule.id":          1,
		"nurse.schedule.name":        1,
		"nurse.schedule.serviceFees": 1,
	}

	medicalLabScientistServiceOpts := options.FindOne().SetProjection(nurseProjection)

	err = serviceColl.FindOne(ctx, nurseServiceFilter, medicalLabScientistServiceOpts).Decode(&service)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(nurse.AppoiynmentResDto{
			Status:  false,
			Message: "Failed to fetch service data: " + err.Error(),
		})
	}

	var nurseServiceData entity.ServiceAndSchedule
	if service.Nurse != nil && len(service.Nurse.Schedule) > 0 {
		for _, nurseService := range service.Nurse.Schedule {
			if nurseService.Id == nurseServiceDataServiceObjID {
				nurseServiceData = nurseService
				break
			}
		}
	}

	longitude, err := strconv.ParseFloat(data.Longitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(nurse.AppoiynmentResDto{
			Status:  false,
			Message: "Invalid longitude format",
		})
	}

	latitude, err := strconv.ParseFloat(data.Latitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(nurse.AppoiynmentResDto{
			Status:  false,
			Message: "Invalid latitude format",
		})
	}

	var remindMeBefore time.Time
	if data.RemindMeBefore != "" {
		remindMeBefore, err = time.Parse(time.DateTime, data.RemindMeBefore)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(nurse.AppoiynmentResDto{
				Status:  false,
				Message: "Failed to parse remindMeBefore date: " + err.Error(),
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(nurse.AppoiynmentResDto{
			Status:  false,
			Message: "remindMeBefore date is mandatory",
		})
	}

	var name string
	var image string

	if service.Nurse != nil {
		name = service.Nurse.Information.Name
		image = service.Nurse.Information.Image
	}

	appointmentData := entity.NurseAppointmentEntity{
		AppointmentDetails: entity.AppointmentDetailsAppointmentEntity{
			From:           fromDate,
			To:             toDate,
			RemindMeBefore: remindMeBefore,
		},
		Service: entity.ServiceAppointmentEntity{
			Id:          nurseServiceData.Id,
			Name:        nurseServiceData.Name,
			ServiceFees: nurseServiceData.ServiceFees,
		},
		Information: entity.NurseInformation{
			Name:  name,
			Image: image,
		},
		Destination: entity.Address{
			Coordinates: []float64{longitude, latitude},
			Add:         data.Address,
			Type:        "Point",
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

	appointmentId := primitive.NewObjectID()
	appointment = entity.AppointmentEntity{
		Id:                   appointmentId,
		Role:                 "healthProfessional",
		FacilityOrProfession: "nurse",
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
		Nurse:             &appointmentData,
		PaymentStatus:     "initiated",
		AppointmentStatus: "pending",
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	session, err := database.GetMongoClient().StartSession()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(nurse.AppoiynmentResDto{
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
			"facilityOrProfession": "nurse",
		}

		update := bson.M{
			"$push": bson.M{
				"nurse.upcommingEvents": bson.M{
					"id":        appointment.Id,
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
		return c.Status(fiber.StatusInternalServerError).JSON(nurse.AppoiynmentResDto{
			Status:  false,
			Message: "Failed to update appointment data: " + err.Error(),
		})
	}

	if service.User.Notification.DeviceToken != "" && service.User.Notification.DeviceType != "" {
		formattedFromDate := fromDate.Format("02 Jan 2006 15:04")
		formattedToDate := toDate.Format("02 Jan 2006 15:04")

		notificationTitle := "New Appointment Scheduled"
		notificationBody := "A new appointment has been scheduled from " + formattedFromDate + " to " + formattedToDate + "."
		notificationData := map[string]string{
			"type":                 "appointment-scheduled",
			"appointmentId":        appointmentId.Hex(),
			"role":                 "healthProfessional",
			"facilityOrProfession": "nurse",
		}

		utils.SendNotificationToUser(service.User.Notification.DeviceToken, service.User.Notification.DeviceType, notificationTitle, notificationBody, notificationData, service.Id, "provider")
	}

	nurseRes := nurse.AppoiynmentResDto{
		Status:  true,
		Message: "Appointment added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(nurseRes)
}
