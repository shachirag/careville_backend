package physiotherapist

import (
	"careville_backend/database"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get appointment by ID
// @Tags provider appointments
// @Description Get appointment by ID
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "appointment ID"
// @Produce json
// @Success 200 {object} services.GetPhysiotherapistAppointmentDetailResDto
// @Router /provider/services/appointment/physiotherapist-appointment/{id} [get]
func GetPhysiotherpistAppointmentByID(c *fiber.Ctx) error {
	var (
		appointmentColl = database.GetCollection("appointment")
	)

	idParam := c.Params("id")
	appointmentID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.GetPhysiotherapistAppointmentDetailResDto{
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
		"customer.age":                              1,
		"facilityOrProfession":                      1,
		"physiotherapist.appointmentDetails.from":   1,
		"physiotherapist.appointmentDetails.to":     1,
		"physiotherapist.pricePaid":                 1,
		"physiotherapist.familyMember.id":           1,
		"physiotherapist.familyMember.name":         1,
		"physiotherapist.familyMember.age":          1,
		"physiotherapist.familyMember.sex":          1,
		"physiotherapist.familyMember.relationship": 1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	var appointment entity.AppointmentEntity
	err = appointmentColl.FindOne(ctx, filter, findOptions).Decode(&appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.GetPhysiotherapistAppointmentDetailResDto{
			Status:  false,
			Message: "Failed to fetch appointment data: " + err.Error(),
		})
	}

	var appointmentFromDate time.Time
	var appointmentToDate time.Time
	var familiyMemberId primitive.ObjectID
	var familiyMemberRelationShip string
	var familiyMemberName string
	var familiyMemberAge string
	var familiyMemberSex string
	var pricePaid float64
	if appointment.Physiotherapist != nil {
		appointmentFromDate = appointment.Physiotherapist.AppointmentDetails.From
		appointmentToDate = appointment.Physiotherapist.AppointmentDetails.To
		familiyMemberId = appointment.Physiotherapist.FamilyMember.ID
		familiyMemberName = appointment.Physiotherapist.FamilyMember.Name
		familiyMemberAge = appointment.Physiotherapist.FamilyMember.Age
		familiyMemberSex = appointment.Physiotherapist.FamilyMember.Sex
		familiyMemberRelationShip = appointment.Physiotherapist.FamilyMember.Relationship
		pricePaid = appointment.Physiotherapist.PricePaid
	}

	expertiseRes := services.GetPhysiotherapistAppointmentDetailResDto{
		Status:  true,
		Message: "Data fetched successfully",
		Data: services.PhysiotherapistAppointmentRes{
			Id: appointment.Id,
			Customer: services.CustomerInformation{
				Id:        appointment.Customer.ID,
				FirstName: appointment.Customer.FirstName,
				LastName:  appointment.Customer.LastName,
				Image:     appointment.Customer.Image,
				PhoneNumber: services.PhoneNumber{
					DialCode:    appointment.Customer.PhoneNumber.DialCode,
					Number:      appointment.Customer.PhoneNumber.Number,
					CountryCode: appointment.Customer.PhoneNumber.CountryCode,
				},
				Age: appointment.Customer.Age,
			},
			FacilityOrProfession: appointment.FacilityOrProfession,
			AppointmentDetails: services.AppointmentDetails{
				AppointmentFromDate: appointmentFromDate,
				AppointmentToDate:   appointmentToDate,
			},
			FamilyMember: services.FamilyMember{
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
