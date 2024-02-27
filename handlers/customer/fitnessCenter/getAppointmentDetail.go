package fitnessCenter

import (
	"careville_backend/database"
	fitnessCenter "careville_backend/dto/customer/fitnessCenter"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get appointment by ID
// @Tags customer appointments
// @Description Get appointment by ID
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "appointment ID"
// @Produce json
// @Success 200 {object} fitnessCenter.GetFitnessCenterAppointmentDetailResDto
// @Router /provider/fitnessCenter/appointment/fitnessCenter-appointment/{id} [get]
func GetFitnessCenterAppointmentByID(c *fiber.Ctx) error {

	var (
		appointmentColl = database.GetCollection("appointment")
		serviceColl     = database.GetCollection("service")
	)

	idParam := c.Params("id")
	appointmentID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fitnessCenter.GetFitnessCenterAppointmentDetailResDto{
			Status:  false,
			Message: "Invalid appointment ID",
		})
	}

	filter := bson.M{"_id": appointmentID}

	projection := bson.M{
		"_id":                1,
		"customer.id":        1,
		"customer.firstName": 1,
		"customer.lastName":  1,
		"customer.image":     1,
		"customer.phoneNumber": bson.M{
			"dialCode":    1,
			"number":      1,
			"countryCode": 1,
		},
		"facilityOrProfession":                    1,
		"doctor.pricePaid":                        1,
		"fitnessCenter.package":                   1,
		"fitnessCenter.trainer.id":                1,
		"fitnessCenter.trainer.name":              1,
		"fitnessCenter.trainer.category":          1,
		"fitnessCenter.trainer.information":       1,
		"fitnessCenter.trainer.price":             1,
		"fitnessCenter.familyMember.id":           1,
		"fitnessCenter.familyMember.name":         1,
		"fitnessCenter.familyMember.age":          1,
		"fitnessCenter.familyMember.sex":          1,
		"fitnessCenter.familyMember.relationship": 1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	var appointment entity.AppointmentEntity
	err = appointmentColl.FindOne(ctx, filter, findOptions).Decode(&appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fitnessCenter.GetFitnessCenterAppointmentDetailResDto{
			Status:  false,
			Message: "Failed to fetch appointment data: " + err.Error(),
		})
	}

	var fitnessCeter entity.ServiceEntity
	reviewFilter := bson.M{"_id": appointment.ServiceID}
	err = serviceColl.FindOne(ctx, reviewFilter, findOptions).Decode(&appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fitnessCenter.GetFitnessCenterAppointmentDetailResDto{
			Status:  false,
			Message: "Failed to fetch average rating: " + err.Error(),
		})
	}

	var avgRating float64
	if fitnessCeter.FitnessCenter != nil {
		avgRating = fitnessCeter.FitnessCenter.Review.AvgRating
	}

	var trainerId primitive.ObjectID
	var trainerName string
	var trainerCategory string
	var trainerInformation string
	var trainerPrice float64
	var gymPackage string
	var subscriptionPrice float64
	var familiyMemberId primitive.ObjectID
	var familiyMemberRelationShip string
	var familiyMemberName string
	var familiyMemberAge string
	var familiyMemberSex string
	var pricePaid float64
	var fitnessCenterImage string
	var fitnessCenterName string
	var fitnessCenterAddress fitnessCenter.Address
	if appointment.FitnessCenter != nil {
		gymPackage = appointment.FitnessCenter.Package
		subscriptionPrice = appointment.FitnessCenter.Invoice.MembershipSubscription
		trainerId = appointment.FitnessCenter.Trainer.ID
		trainerName = appointment.FitnessCenter.Trainer.Name
		trainerCategory = appointment.FitnessCenter.Trainer.Category
		trainerInformation = appointment.FitnessCenter.Trainer.Information
		trainerPrice = appointment.FitnessCenter.Trainer.Price
		familiyMemberId = appointment.FitnessCenter.FamilyMember.ID
		familiyMemberName = appointment.FitnessCenter.FamilyMember.Name
		familiyMemberAge = appointment.FitnessCenter.FamilyMember.Age
		familiyMemberSex = appointment.FitnessCenter.FamilyMember.Sex
		familiyMemberRelationShip = appointment.FitnessCenter.FamilyMember.Relationship
		pricePaid = appointment.FitnessCenter.Invoice.TotalAmountPaid
		fitnessCenterName = appointment.FitnessCenter.Information.Name
		fitnessCenterImage = appointment.FitnessCenter.Information.Image
		fitnessCenterAddress = fitnessCenter.Address(appointment.FitnessCenter.Information.Address)
	}

	expertiseRes := fitnessCenter.GetFitnessCenterAppointmentDetailResDto{
		Status:  true,
		Message: "Data fetched successfully",
		Data: fitnessCenter.FitnessCenterAppointmentRes{
			Id: appointment.Id,
			Customer: fitnessCenter.CustomerInformation{
				Id:        appointment.Customer.ID,
				FirstName: appointment.Customer.FirstName,
				LastName:  appointment.Customer.LastName,
				Image:     appointment.Customer.Image,
				PhoneNumber: fitnessCenter.PhoneNumber{
					DialCode:    appointment.Customer.PhoneNumber.DialCode,
					Number:      appointment.Customer.PhoneNumber.Number,
					CountryCode: appointment.Customer.PhoneNumber.CountryCode,
				},
			},
			FitnessCenterInformation: fitnessCenter.FitnessCenterInformation{
				Id:        appointment.ServiceID,
				Name:      fitnessCenterName,
				Image:     fitnessCenterImage,
				Address:   fitnessCenterAddress,
				AvgRating: avgRating,
			},
			Subscription: fitnessCenter.SubscriptionData{
				Package: gymPackage,
				Price:   subscriptionPrice,
			},
			FacilityOrProfession: appointment.FacilityOrProfession,
			TrainerInformation: fitnessCenter.TrainerInformation{
				Id:          trainerId,
				Name:        trainerName,
				Category:    trainerCategory,
				Price:       trainerPrice,
				Information: trainerInformation,
			},
			FamilyMember: fitnessCenter.FamilyMember{
				Id:           familiyMemberId,
				Name:         familiyMemberName,
				Age:          familiyMemberAge,
				Sex:          familiyMemberSex,
				RelationShip: familiyMemberRelationShip,
			},
			PricePaid: pricePaid,
		},
	}

	return c.Status(fiber.StatusOK).JSON(expertiseRes)
}
