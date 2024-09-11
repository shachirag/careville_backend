package nurse

import (
	"careville_backend/database"
	nurse "careville_backend/dto/customer/nurse"
	"careville_backend/entity"
	"time"

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
// @Success 200 {object} nurse.GetNurseAppointmentDetailResDto
// @Router /customer/healthProfessional/appointment/nurse-appointment/{id} [get]
func GetNurseAppointmentByID(c *fiber.Ctx) error {
	var (
		appointmentColl = database.GetCollection("appointment")
	)

	idParam := c.Params("id")
	appointmentID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(nurse.GetNurseAppointmentDetailResDto{
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
		"age":                             1,
		"facilityOrProfession":            1,
		"nurse.appointmentDetails.from":   1,
		"nurse.appointmentDetails.to":     1,
		"nurse.pricePaid":                 1,
		"nurse.familyMember.id":           1,
		"nurse.familyMember.name":         1,
		"nurse.familyMember.age":          1,
		"nurse.familyMember.sex":          1,
		"nurse.familyMember.relationship": 1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	var appointment entity.AppointmentEntity
	err = appointmentColl.FindOne(ctx, filter, findOptions).Decode(&appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(nurse.GetNurseAppointmentDetailResDto{
			Status:  false,
			Message: "Failed to fetch appointment data: " + err.Error(),
		})
	}

	var nurse1 entity.ServiceEntity
	reviewFilter := bson.M{"_id": appointment.ServiceID}
	projection = bson.M{
		"nurse.review.avgRating": 1,
	}

	reviewFindOptions := options.FindOne().SetProjection(projection)
	err = database.GetCollection("service").FindOne(ctx, reviewFilter, reviewFindOptions).Decode(&appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(nurse.GetNurseAppointmentDetailResDto{
			Status:  false,
			Message: "Failed to fetch average rating: " + err.Error(),
		})
	}

	var avgRating float64
	if nurse1.Nurse != nil {
		avgRating = nurse1.Nurse.Review.AvgRating
	}

	var appointmentFromDate time.Time
	var appointmentToDate time.Time
	var familiyMemberId primitive.ObjectID
	var familiyMemberRelationShip string
	var familiyMemberName string
	var familiyMemberAge string
	var familiyMemberSex string
	var pricePaid float64
	if appointment.Nurse != nil {
		appointmentFromDate = appointment.Nurse.AppointmentDetails.From
		appointmentToDate = appointment.Nurse.AppointmentDetails.To
		familiyMemberId = appointment.Nurse.FamilyMember.ID
		familiyMemberName = appointment.Nurse.FamilyMember.Name
		familiyMemberAge = appointment.Nurse.FamilyMember.Age
		familiyMemberSex = appointment.Nurse.FamilyMember.Sex
		familiyMemberRelationShip = appointment.Nurse.FamilyMember.Relationship
		pricePaid = appointment.Nurse.PricePaid
	}

	expertiseRes := nurse.GetNurseAppointmentDetailResDto{
		Status:  true,
		Message: "Data fetched successfully",
		Data: nurse.NurseAppointmentRes{
			Id: appointment.Id,
			Customer: nurse.CustomerInformation{
				Id:        appointment.Customer.ID,
				FirstName: appointment.Customer.FirstName,
				LastName:  appointment.Customer.LastName,
				Image:     appointment.Customer.Image,
				PhoneNumber: nurse.PhoneNumber{
					DialCode:    appointment.Customer.PhoneNumber.DialCode,
					Number:      appointment.Customer.PhoneNumber.Number,
					CountryCode: appointment.Customer.PhoneNumber.CountryCode,
				},
				Age: appointment.Customer.Age,
			},
			AvgRating:            avgRating,
			FacilityOrProfession: appointment.FacilityOrProfession,
			AppointmentDetails: nurse.AppointmentDetails{
				AppointmentFromDate: appointmentFromDate,
				AppointmentToDate:   appointmentToDate,
			},
			FamilyMember: nurse.FamilyMember{
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
