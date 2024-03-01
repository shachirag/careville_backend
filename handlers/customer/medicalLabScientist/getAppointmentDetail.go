package medicalLabScientist

import (
	"careville_backend/database"
	medicalLabScientist "careville_backend/dto/customer/medicalLabScientist"
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
// @Success 200 {object} medicalLabScientist.GetMedicalLabScientistAppointmentDetailResDto
// @Router /customer/healthProfessional/appointment/medicalLabScientist-appointment/{id} [get]
func GetMedicalLabScientistAppointmentByID(c *fiber.Ctx) error {
	var (
		appointmentColl = database.GetCollection("appointment")
	)

	idParam := c.Params("id")
	appointmentID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(medicalLabScientist.GetMedicalLabScientistAppointmentDetailResDto{
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
		"facilityOrProfession":                          1,
		"medicalLabScientist.appointmentDetails.from":   1,
		"medicalLabScientist.appointmentDetails.to":     1,
		"medicalLabScientist.pricePaid":                 1,
		"medicalLabScientist.familyMember.id":           1,
		"medicalLabScientist.familyMember.name":         1,
		"medicalLabScientist.familyMember.age":          1,
		"medicalLabScientist.familyMember.sex":          1,
		"medicalLabScientist.familyMember.relationship": 1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	var appointment entity.AppointmentEntity
	err = appointmentColl.FindOne(ctx, filter, findOptions).Decode(&appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.GetMedicalLabScientistAppointmentDetailResDto{
			Status:  false,
			Message: "Failed to fetch appointment data: " + err.Error(),
		})
	}

	var medicalLabScientist1 entity.ServiceEntity
	reviewFilter := bson.M{"_id": appointment.ServiceID}
	projection = bson.M{
		"medicalLabScientist.review.avgRating": 1,
	}

	reviewFindOptions := options.FindOne().SetProjection(projection)
	err = database.GetCollection("service").FindOne(ctx, reviewFilter, reviewFindOptions).Decode(&appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.GetMedicalLabScientistAppointmentDetailResDto{
			Status:  false,
			Message: "Failed to fetch average rating: " + err.Error(),
		})
	}

	var avgRating float64
	if medicalLabScientist1.MedicalLabScientist != nil {
		avgRating = medicalLabScientist1.MedicalLabScientist.Review.AvgRating
	}

	var appointmentFromDate time.Time
	var appointmentToDate time.Time
	var familiyMemberId primitive.ObjectID
	var familiyMemberRelationShip string
	var familiyMemberName string
	var familiyMemberAge string
	var familiyMemberSex string
	var pricePaid float64
	if appointment.MedicalLabScientist != nil {
		appointmentFromDate = appointment.MedicalLabScientist.AppointmentDetails.From
		appointmentToDate = appointment.MedicalLabScientist.AppointmentDetails.To
		familiyMemberId = appointment.MedicalLabScientist.FamilyMember.ID
		familiyMemberName = appointment.MedicalLabScientist.FamilyMember.Name
		familiyMemberAge = appointment.MedicalLabScientist.FamilyMember.Age
		familiyMemberSex = appointment.MedicalLabScientist.FamilyMember.Sex
		familiyMemberRelationShip = appointment.MedicalLabScientist.FamilyMember.Relationship
		pricePaid = appointment.MedicalLabScientist.PricePaid
	}

	expertiseRes := medicalLabScientist.GetMedicalLabScientistAppointmentDetailResDto{
		Status:  true,
		Message: "Data fetched successfully",
		Data: medicalLabScientist.MedicallabScientistAppointmentRes{
			Id: appointment.Id,
			Customer: medicalLabScientist.CustomerInformation{
				Id:        appointment.Customer.ID,
				FirstName: appointment.Customer.FirstName,
				LastName:  appointment.Customer.LastName,
				Image:     appointment.Customer.Image,
				PhoneNumber: medicalLabScientist.PhoneNumber{
					DialCode:    appointment.Customer.PhoneNumber.DialCode,
					Number:      appointment.Customer.PhoneNumber.Number,
					CountryCode: appointment.Customer.PhoneNumber.CountryCode,
				},
			},
			AvgRating:            avgRating,
			FacilityOrProfession: appointment.FacilityOrProfession,
			AppointmentDetails: medicalLabScientist.AppointmentDetails{
				AppointmentFromDate: appointmentFromDate,
				AppointmentToDate:   appointmentToDate,
			},
			FamilyMember: medicalLabScientist.FamilyMember{
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
