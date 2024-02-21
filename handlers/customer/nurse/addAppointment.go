package nurse

import (
	"strconv"
	"time"

	"careville_backend/database"
	customerMiddleware "careville_backend/dto/customer/middleware"
	"careville_backend/dto/customer/nurse"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
// @Success 200 {object}  nurse.NurseAppointmentResDto
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
		return c.Status(500).JSON(nurse.NurseAppointmentResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	var familyObjectID primitive.ObjectID

	if data.FamillyMemberId != nil && *data.FamillyMemberId != "" {

		familyObjectID, err = primitive.ObjectIDFromHex(*data.FamillyMemberId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(nurse.NurseAppointmentResDto{
				Status:  false,
				Message: "Invalid ID format",
			})
		}
	}

	serviceId := c.Query("serviceId")

	if serviceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(nurse.NurseAppointmentResDto{
			Status:  false,
			Message: "service Id is mandatory",
		})
	}

	serviceObjectID, err := primitive.ObjectIDFromHex(serviceId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(nurse.NurseAppointmentResDto{
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
			return c.Status(fiber.StatusInternalServerError).JSON(nurse.NurseAppointmentResDto{
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
		return c.Status(fiber.StatusInternalServerError).JSON(nurse.NurseAppointmentResDto{
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
		"nurse.information.name":  1,
		"nurse.information.image": 1,
	}

	serviceOpts := options.FindOne().SetProjection(serviceProjection)

	var service entity.ServiceEntity
	err = serviceColl.FindOne(ctx, serviceFilter, serviceOpts).Decode(&service)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(nurse.NurseAppointmentResDto{
			Status:  false,
			Message: "Failed to fetch nurse data: " + err.Error(),
		})
	}

	if service.Nurse == nil {
		return c.Status(fiber.StatusNotFound).JSON(nurse.NurseAppointmentResDto{
			Status:  false,
			Message: "nurse data not found",
		})
	}

	longitude, err := strconv.ParseFloat(data.Longitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(nurse.NurseAppointmentResDto{
			Status:  false,
			Message: "Invalid longitude format",
		})
	}

	latitude, err := strconv.ParseFloat(data.Latitude, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(nurse.NurseAppointmentResDto{
			Status:  false,
			Message: "Invalid latitude format",
		})
	}

	var fromDate time.Time
	if data.FromDate != "" {
		fromDate, err = time.Parse(time.DateTime, data.FromDate)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(nurse.NurseAppointmentResDto{
				Status:  false,
				Message: "Failed to parse fromDate date: " + err.Error(),
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(nurse.NurseAppointmentResDto{
			Status:  false,
			Message: "fromDate is mandatory",
		})
	}

	var toDate time.Time
	if data.ToDate != "" {
		toDate, err = time.Parse(time.DateTime, data.ToDate)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(nurse.NurseAppointmentResDto{
				Status:  false,
				Message: "Failed to parse toDate date: " + err.Error(),
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(nurse.NurseAppointmentResDto{
			Status:  false,
			Message: "toDate date is mandatory",
		})
	}

	var remindMeBefore time.Time
	if data.RemindMeBefore != "" {
		remindMeBefore, err = time.Parse(time.DateTime, data.RemindMeBefore)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(nurse.NurseAppointmentResDto{
				Status:  false,
				Message: "Failed to parse remindMeBefore date: " + err.Error(),
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(nurse.NurseAppointmentResDto{
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

	appointment = entity.AppointmentEntity{
		Id:                   primitive.NewObjectID(),
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
		},
		Nurse:             &appointmentData,
		PaymentStatus:     "initiated",
		AppointmentStatus: "pending",
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	_, err = appointmentColl.InsertOne(ctx, appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(nurse.NurseAppointmentResDto{
			Status:  false,
			Message: "Failed to insert nurse appointment data into MongoDB: " + err.Error(),
		})
	}

	nurseRes := nurse.NurseAppointmentResDto{
		Status:  true,
		Message: "Appointment added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(nurseRes)
}
