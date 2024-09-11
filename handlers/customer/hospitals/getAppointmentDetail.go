package hospitals

import (
	"careville_backend/database"
	hospitals "careville_backend/dto/customer/hospitals"
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
// @Success 200 {object} hospitals.GetHospitalAppointmentDetailResDto
// @Router /customer/healthFacility/appointment/hospital-appointment/{id} [get]
func GetHospitalAppointmentByID(c *fiber.Ctx) error {
	var (
		appointmentColl = database.GetCollection("appointment")
		// serviceColl     = database.GetCollection("service")
	)

	idParam := c.Params("id")
	appointmentID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(hospitals.GetHospitalAppointmentDetailResDto{
			Status:  false,
			Message: "Invalid appointment ID",
		})
	}

	filter := bson.M{"_id": appointmentID}

	projection := bson.M{
		"_id":                1,
		"serviceId":          1,
		"customer.id":        1,
		"customer.firstName": 1,
		"customer.lastName":  1,
		"customer.image":     1,
		"customer.phoneNumber": bson.M{
			"dialCode":    1,
			"number":      1,
			"countryCode": 1,
		},
		"age":                          1,
		"facilityOrProfession":         1,
		"hospClinic.information.name":  1,
		"hospClinic.information.image": 1,
		"hospClinic.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
		"hospital.appointmentDetails.from":   1,
		"hospital.appointmentDetails.to":     1,
		"hospital.pricePaid":                 1,
		"hospital.familyMember.id":           1,
		"hospital.familyMember.name":         1,
		"hospital.familyMember.age":          1,
		"hospital.familyMember.sex":          1,
		"hospital.familyMember.relationship": 1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	var appointment entity.AppointmentEntity
	err = appointmentColl.FindOne(ctx, filter, findOptions).Decode(&appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(hospitals.GetHospitalAppointmentDetailResDto{
			Status:  false,
			Message: "Failed to fetch appointment data: " + err.Error(),
		})
	}

	var hospital entity.ServiceEntity
	reviewFilter := bson.M{"_id": appointment.ServiceID}
	projection = bson.M{
		"hospClinic.review.avgRating": 1,
	}

	reviewFindOptions := options.FindOne().SetProjection(projection)
	err = database.GetCollection("service").FindOne(ctx, reviewFilter, reviewFindOptions).Decode(&appointment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(hospitals.GetHospitalAppointmentDetailResDto{
			Status:  false,
			Message: "Failed to fetch average rating: " + err.Error(),
		})
	}

	var avgRating float64
	if hospital.HospClinic != nil {
		avgRating = hospital.HospClinic.Review.AvgRating
	}

	var appointmentFromDate time.Time
	var appointmentToDate time.Time
	var familiyMemberId primitive.ObjectID
	var familiyMemberRelationShip string
	var familiyMemberName string
	var familiyMemberAge string
	var familiyMemberSex string
	var pricePaid float64
	var hospitalImage string
	var hospitalName string
	var hospitalAddress hospitals.Address
	if appointment.HospitalClinic != nil {
		appointmentFromDate = appointment.HospitalClinic.AppointmentDetails.From
		appointmentToDate = appointment.HospitalClinic.AppointmentDetails.To
		familiyMemberId = appointment.HospitalClinic.FamilyMember.ID
		familiyMemberName = appointment.HospitalClinic.FamilyMember.Name
		familiyMemberAge = appointment.HospitalClinic.FamilyMember.Age
		familiyMemberSex = appointment.HospitalClinic.FamilyMember.Sex
		familiyMemberRelationShip = appointment.HospitalClinic.FamilyMember.Relationship
		pricePaid = appointment.HospitalClinic.PricePaid
		hospitalName = appointment.HospitalClinic.Information.Name
		hospitalImage = appointment.HospitalClinic.Information.Image
		hospitalAddress = hospitals.Address(appointment.HospitalClinic.Information.Address)
	}

	expertiseRes := hospitals.GetHospitalAppointmentDetailResDto{
		Status:  true,
		Message: "Data fetched successfully",
		Data: hospitals.HospitalAppointmentRes{
			Id: appointment.Id,
			Customer: hospitals.CustomerInformation{
				Id:        appointment.Customer.ID,
				FirstName: appointment.Customer.FirstName,
				LastName:  appointment.Customer.LastName,
				Image:     appointment.Customer.Image,
				PhoneNumber: hospitals.PhoneNumber{
					DialCode:    appointment.Customer.PhoneNumber.DialCode,
					Number:      appointment.Customer.PhoneNumber.Number,
					CountryCode: appointment.Customer.PhoneNumber.CountryCode,
				},
				Age: appointment.Customer.Age,
			},
			HospitalInformation: hospitals.HospitalInformation{
				Id:        appointment.ServiceID,
				Name:      hospitalName,
				Image:     hospitalImage,
				Address:   hospitalAddress,
				AvgRating: avgRating,
			},
			FacilityOrProfession: appointment.FacilityOrProfession,
			AppointmentDetails: hospitals.AppointmentDetails{
				AppointmentFromDate: appointmentFromDate,
				AppointmentToDate:   appointmentToDate,
			},
			FamilyMember: hospitals.FamilyMember{
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
