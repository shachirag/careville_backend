package laboratory

import (
	"time"

	"careville_backend/database"
	laboratory "careville_backend/dto/customer/laboratories"
	customerMiddleware "careville_backend/dto/customer/middleware"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Add appointment
// @Tags customer laboratory
// @Description Add appointment
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param serviceId query string true "service ID"
// @Param  customer body laboratory.LaboratoryAppointmentReqDto true "add laboratory appointment"
// @Produce json
// @Success 200 {object} laboratory.LaboratoryAppointmentResDto
// @Router /customer/healthFacility/add-laboratory-appointment [post]
func AddLaboratoryAppointment(c *fiber.Ctx) error {

	var (
		appointmentColl = database.GetCollection("appointment")
		serviceColl     = database.GetCollection("service")
		customerColl    = database.GetCollection("customer")
		data            laboratory.LaboratoryAppointmentReqDto
		appointment     entity.AppointmentEntity
		service         entity.ServiceEntity
	)

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(laboratory.LaboratoryAppointmentResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	familyObjectID, err := primitive.ObjectIDFromHex(data.FamillyMemberId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(laboratory.LaboratoryAppointmentResDto{
			Status:  false,
			Message: "Invalid ID format",
		})
	}

	serviceId := c.Query("serviceId")

	if serviceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(laboratory.LaboratoryAppointmentResDto{
			Status:  false,
			Message: "service Id is mandatory",
		})
	}

	serviceObjectID, err := primitive.ObjectIDFromHex(serviceId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(laboratory.LaboratoryAppointmentResDto{
			Status:  false,
			Message: "Invalid ID format",
		})
	}

	investigationObjID, err := primitive.ObjectIDFromHex(data.InvestigationId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(laboratory.LaboratoryAppointmentResDto{
			Status:  false,
			Message: "Invalid ID format",
		})
	}

	investigationFilter := bson.M{
		"_id": serviceObjectID,
		"laboratory.investigations": bson.M{
			"$elemMatch": bson.M{
				"id": investigationObjID,
			},
		},
	}

	investigationProjection := bson.M{
		"laboratory.investigations.id":          1,
		"laboratory.investigations.name":        1,
		"laboratory.investigations.type":        1,
		"laboratory.investigations.information": 1,
		"laboratory.investigations.price":       1,
	}

	investigationOpts := options.FindOne().SetProjection(investigationProjection)

	err = serviceColl.FindOne(ctx, investigationFilter, investigationOpts).Decode(&service)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(laboratory.LaboratoryAppointmentResDto{
			Status:  false,
			Message: "Failed to fetch investigation data: " + err.Error(),
		})
	}

	if service.Laboratory == nil {
		return c.Status(fiber.StatusNotFound).JSON(laboratory.LaboratoryAppointmentResDto{
			Status:  false,
			Message: "Laboratory data not found",
		})
	}

	var investigationData entity.Investigations
	if service.Laboratory != nil && len(service.Laboratory.Investigations) > 0 {
		for _, investiagtion := range service.Laboratory.Investigations {
			if investiagtion.Id == investigationObjID {
				investigationData = investiagtion
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
		return c.Status(fiber.StatusInternalServerError).JSON(laboratory.LaboratoryAppointmentResDto{
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
		return c.Status(fiber.StatusInternalServerError).JSON(laboratory.LaboratoryAppointmentResDto{
			Status:  false,
			Message: "Failed to fetch customer data: " + err.Error(),
		})
	}

	var appointmentDate time.Time
	if data.AppointmentDate != "" {
		appointmentDate, err = time.Parse(time.RFC3339, data.AppointmentDate)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(laboratory.LaboratoryAppointmentResDto{
				Status:  false,
				Message: "Failed to parse appointment date: " + err.Error(),
			})
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(laboratory.LaboratoryAppointmentResDto{
			Status:  false,
			Message: "Appointment date is mandatory",
		})
	}

	appointmentData := entity.LaboratoryAppointmentEntity{
		Investigation: entity.InvestigationAppointmentEntity{
			ID:          investigationObjID,
			Name:        investigationData.Name,
			Information: investigationData.Information,
			Type:        investigationData.Type,
			Price:       investigationData.Price,
		},
		FamilyMember: entity.FamilyMemberAppointmentEntity{
			ID:           familyObjectID,
			Name:         familyData.Name,
			Age:          familyData.Age,
			Sex:          familyData.Sex,
			Relationship: familyData.RelationShip,
		},
		AppointmentDetails: entity.LaboratoryAppointmentDetailsAppointmentEntity{
			Date: appointmentDate,
		},
		FamilyType: data.FamilyType,
		PricePaid:  data.PricePaid,
	}

	appointment = entity.AppointmentEntity{
		Id:                   primitive.NewObjectID(),
		Role:                 "healthFacility",
		FacilityOrProfession: "laboratory",
		ServiceID:            serviceObjectID,
		Customer: entity.CustomerAppointmentEntity{
			ID:          customerMiddlewareData.CustomerId,
			FirstName:   customer.FirstName,
			LastName:    customer.LastName,
			Image:       customer.Image,
			Email:       customer.Email,
			PhoneNumber: customer.PhoneNumber,
		},
		Laboratory:        &appointmentData,
		PaymentStatus:     "initiated",
		AppointmentStatus: "pending",
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	_, err = appointmentColl.InsertOne(ctx, appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(laboratory.LaboratoryAppointmentResDto{
			Status:  false,
			Message: "Failed to insert laboratory appointment data into MongoDB: " + err.Error(),
		})
	}

	hospClinicRes := laboratory.LaboratoryAppointmentResDto{
		Status:  true,
		Message: "Appointment added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(hospClinicRes)
}
