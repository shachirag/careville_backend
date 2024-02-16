package medicalLabScientist

import (
	"time"

	"careville_backend/database"
	"careville_backend/dto/customer/medicalLabScientist"
	customerMiddleware "careville_backend/dto/customer/middleware"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	var fromDate time.Time
	if data.FromDate != "" {
		fromDate, err = time.Parse(time.RFC3339, data.FromDate)
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
		toDate, err = time.Parse(time.RFC3339, data.ToDate)
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

	var remindMeBefore time.Time
	if data.RemindMeBefore != "" {
		remindMeBefore, err = time.Parse(time.RFC3339, data.RemindMeBefore)
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

	appointment = entity.AppointmentEntity{
		Id:                   primitive.NewObjectID(),
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
		},
		MedicalLabScientist: &appointmentData,
		PaymentStatus:       "initiated",
		AppointmentStatus:   "pending",
		CreatedAt:           time.Now().UTC(),
		UpdatedAt:           time.Now().UTC(),
	}

	_, err = appointmentColl.InsertOne(ctx, appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.MedicalLabScientistAppointmentResDto{
			Status:  false,
			Message: "Failed to insert medicalLabScientist appointment data into MongoDB: " + err.Error(),
		})
	}

	medicalLabScientistRes := medicalLabScientist.MedicalLabScientistAppointmentResDto{
		Status:  true,
		Message: "Appointment added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(medicalLabScientistRes)
}
