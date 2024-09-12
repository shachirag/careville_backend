package doctorProfession

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
// @Success 200 {object} services.GetDoctorProfessionAppointmentDetailResDto
// @Router /provider/services/appointment/doctor-appointment/{id} [get]
func GetDoctorprofessionAppointmentByID(c *fiber.Ctx) error {
	var (
		appointmentColl = database.GetCollection("appointment")
	)

	idParam := c.Params("id")
	appointmentID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(services.GetDoctorProfessionAppointmentDetailResDto{
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
		"customer.age":                     1,
		"facilityOrProfession":             1,
		"doctor.appointmentDetails.from":   1,
		"doctor.appointmentDetails.to":     1,
		"doctor.pricePaid":                 1,
		"doctor.familyMember.id":           1,
		"doctor.familyMember.name":         1,
		"doctor.familyMember.age":          1,
		"doctor.familyMember.sex":          1,
		"doctor.familyMember.relationship": 1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	var appointment entity.AppointmentEntity
	err = appointmentColl.FindOne(ctx, filter, findOptions).Decode(&appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.GetDoctorProfessionAppointmentDetailResDto{
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
	if appointment.Doctor != nil {
		appointmentFromDate = appointment.Doctor.AppointmentDetails.From
		appointmentToDate = appointment.Doctor.AppointmentDetails.To
		familiyMemberId = appointment.Doctor.FamilyMember.ID
		familiyMemberName = appointment.Doctor.FamilyMember.Name
		familiyMemberAge = appointment.Doctor.FamilyMember.Age
		familiyMemberSex = appointment.Doctor.FamilyMember.Sex
		familiyMemberRelationShip = appointment.Doctor.FamilyMember.Relationship
		pricePaid = appointment.Doctor.PricePaid
	}

	expertiseRes := services.GetDoctorProfessionAppointmentDetailResDto{
		Status:  true,
		Message: "Data fetched successfully",
		Data: services.DoctorProfessionAppointmentRes{
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
