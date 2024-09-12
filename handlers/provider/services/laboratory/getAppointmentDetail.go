package laboratory

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
// @Success 200 {object} services.GetLaboratoryAppointmentDetailResDto
// @Router /provider/services/appointment/laboratory-appointment/{id} [get]
func GetLaboratoryAppointmentByID(c *fiber.Ctx) error {
	var (
		appointmentColl = database.GetCollection("appointment")
	)

	idParam := c.Params("id")
	appointmentID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.GetLaboratoryAppointmentDetailResDto{
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
		"customer.age":                         1,
		"facilityOrProfession":                 1,
		"laboratory.appointmentDetails.date":   1,
		"laboratory.investigation.id":          1,
		"laboratory.investigation.name":        1,
		"laboratory.investigation.information": 1,
		"laboratory.investigation.type":        1,
		"laboratory.investigation.price":       1,
		"laboratory.pricePaid":                 1,
		"laboratory.familyMember.id":           1,
		"laboratory.familyMember.name":         1,
		"laboratory.familyMember.age":          1,
		"laboratory.familyMember.sex":          1,
		"laboratory.familyMember.relationship": 1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	var appointment entity.AppointmentEntity
	err = appointmentColl.FindOne(ctx, filter, findOptions).Decode(&appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.GetLaboratoryAppointmentDetailResDto{
			Status:  false,
			Message: "Failed to fetch appointment data: " + err.Error(),
		})
	}

	var appointmentDate time.Time
	var familiyMemberId primitive.ObjectID
	var familiyMemberRelationShip string
	var familiyMemberName string
	var familiyMemberAge string
	var familiyMemberSex string
	var investigationId primitive.ObjectID
	var investigationName string
	var investigationInformation string
	var investigationType string
	var investigationPrice float64
	var pricePaid float64
	if appointment.Laboratory != nil {
		appointmentDate = appointment.Laboratory.AppointmentDetails.Date
		familiyMemberId = appointment.Laboratory.FamilyMember.ID
		familiyMemberName = appointment.Laboratory.FamilyMember.Name
		familiyMemberAge = appointment.Laboratory.FamilyMember.Age
		familiyMemberSex = appointment.Laboratory.FamilyMember.Sex
		familiyMemberRelationShip = appointment.Laboratory.FamilyMember.Relationship
		investigationId = appointment.Laboratory.Investigation.ID
		investigationName = appointment.Laboratory.Investigation.Name
		investigationInformation = appointment.Laboratory.Investigation.Information
		investigationType = appointment.Laboratory.Investigation.Type
		investigationPrice = appointment.Laboratory.Investigation.Price
		pricePaid = appointment.Laboratory.PricePaid
	}

	expertiseRes := services.GetLaboratoryAppointmentDetailResDto{
		Status:  true,
		Message: "Data fetched successfully",
		Data: services.LaboratoryAppointmentRes{
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
			AppointmentDetails: services.AppointmentData{
				AppointmentDate: appointmentDate,
			},
			Investigation: services.Investigation{
				ID:          investigationId,
				Name:        investigationName,
				Information: investigationInformation,
				Type:        investigationType,
				Price:       investigationPrice,
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
