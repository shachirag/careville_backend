package fitnessCenter

import (
	"time"

	"careville_backend/database"
	"careville_backend/dto/customer/fitnessCenter"
	customerMiddleware "careville_backend/dto/customer/middleware"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Add appointment
// @Tags customer fitnessCenter
// @Description Add appointment
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param serviceId query string true "service ID"
// @Param  customer body fitnessCenter.FitnessCenterAppointmentReqDto true "add HospitalClinic"
// @Produce json
// @Success 200 {object} fitnessCenter.FitnessCenterAppointmentResDto
// @Router /customer/healthFacility/add-fitnessCenter-appointment [post]
func AddFitnessCenterAppointment(c *fiber.Ctx) error {

	var (
		appointmentColl = database.GetCollection("appointment")
		serviceColl     = database.GetCollection("service")
		customerColl    = database.GetCollection("customer")
		data            fitnessCenter.FitnessCenterAppointmentReqDto
		appointment     entity.AppointmentEntity
		service         entity.ServiceEntity
	)

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(fitnessCenter.FitnessCenterAppointmentResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	familyObjectID, err := primitive.ObjectIDFromHex(data.FamillyMemberId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fitnessCenter.FitnessCenterAppointmentResDto{
			Status:  false,
			Message: "Invalid ID format",
		})
	}

	serviceId := c.Query("serviceId")

	if serviceId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fitnessCenter.FitnessCenterAppointmentResDto{
			Status:  false,
			Message: "service Id is mandatory",
		})
	}

	serviceObjectID, err := primitive.ObjectIDFromHex(serviceId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fitnessCenter.FitnessCenterAppointmentResDto{
			Status:  false,
			Message: "Invalid ID format",
		})
	}

	trainerObjID, err := primitive.ObjectIDFromHex(data.TrainerId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fitnessCenter.FitnessCenterAppointmentResDto{
			Status:  false,
			Message: "Invalid ID format",
		})
	}

	trainerFilter := bson.M{
		"_id": serviceObjectID,
		"fitnessCenter.trainers": bson.M{
			"$elemMatch": bson.M{
				"id": trainerObjID,
			},
		},
	}

	trainerProjection := bson.M{
		"fitnessCenter.trainers.id":          1,
		"fitnessCenter.trainers.name":        1,
		"fitnessCenter.trainers.price":       1,
		"fitnessCenter.trainers.information": 1,
		"fitnessCenter.trainers.category":    1,
	}

	trainerOpts := options.FindOne().SetProjection(trainerProjection)

	err = serviceColl.FindOne(ctx, trainerFilter, trainerOpts).Decode(&service)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fitnessCenter.FitnessCenterAppointmentResDto{
			Status:  false,
			Message: "Failed to fetch trainer data: " + err.Error(),
		})
	}

	if service.FitnessCenter == nil {
		return c.Status(fiber.StatusNotFound).JSON(fitnessCenter.FitnessCenterAppointmentResDto{
			Status:  false,
			Message: "Fitness Center data not found",
		})
	}

	var trainerData entity.Trainers
	if service.FitnessCenter != nil && len(service.FitnessCenter.Trainers) > 0 {
		for _, center := range service.FitnessCenter.Trainers {
			if center.Id == trainerObjID {
				trainerData = center
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
		return c.Status(fiber.StatusInternalServerError).JSON(fitnessCenter.FitnessCenterAppointmentResDto{
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
		return c.Status(fiber.StatusInternalServerError).JSON(fitnessCenter.FitnessCenterAppointmentResDto{
			Status:  false,
			Message: "Failed to fetch customer data: " + err.Error(),
		})
	}

	appointmentData := entity.FitnessCenterAppointmentEntity{
		Package: data.Package,
		Trainer: entity.TrainerAppointmentEntity{
			ID:          trainerObjID,
			Category:    trainerData.Category,
			Name:        trainerData.Name,
			Information: trainerData.Information,
			Price:       trainerData.Price,
		},
		FamilyMember: entity.FamilyMemberAppointmentEntity{
			ID:           familyObjectID,
			Name:         familyData.Name,
			Age:          familyData.Age,
			Sex:          familyData.Sex,
			Relationship: familyData.RelationShip,
		},
		FamilyType: data.FamilyType,
		Invoice: entity.Invoice{
			TrainerFees:            trainerData.Price,
			MembershipSubscription: data.MembershipSubscription,
			TotalAmountPaid:        trainerData.Price + data.MembershipSubscription,
		},
	}

	appointment = entity.AppointmentEntity{
		Id:                   primitive.NewObjectID(),
		Role:                 "healthFacility",
		FacilityOrProfession: "fitnessCenter",
		ServiceID:            serviceObjectID,
		Customer: entity.CustomerAppointmentEntity{
			ID:          customerMiddlewareData.CustomerId,
			FirstName:   customer.FirstName,
			LastName:    customer.LastName,
			Image:       customer.Image,
			Email:       customer.Email,
			PhoneNumber: customer.PhoneNumber,
		},
		FitnessCenter:     &appointmentData,
		PaymentStatus:     "initiated",
		AppointmentStatus: "pending",
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	_, err = appointmentColl.InsertOne(ctx, appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fitnessCenter.FitnessCenterAppointmentResDto{
			Status:  false,
			Message: "Failed to insert fitnessCenter appointment data into MongoDB: " + err.Error(),
		})
	}

	fitnessCenterRes := fitnessCenter.FitnessCenterAppointmentResDto{
		Status:  true,
		Message: "Appointment added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(fitnessCenterRes)
}
