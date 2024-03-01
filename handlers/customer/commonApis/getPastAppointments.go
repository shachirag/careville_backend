package common

import (
	"careville_backend/database"
	common "careville_backend/dto/customer/commonApis"
	"careville_backend/entity"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Fetch past appointments
// @Description Fetch past appointments
// @Tags customer appointments
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param customerId query string true "customer ID"
// @Param page query int false "Page no. to fetch the products for 1"
// @Param perPage query int false "Limit of products to fetch is 15"
// @Produce json
// @Success 200 {object} common.GetPastAppointmentsPaginationRes
// @Router /customer/healthFacility/appointment/past-appointments [get]
func FetchPastAppointmentsWithPagination(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "15"))

	appointmentColl := database.GetCollection("appointment")

	customerId := c.Query("customerId")
	customerObjID, err := primitive.ObjectIDFromHex(customerId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.GetPastAppointmentsPaginationRes{
			Status:  false,
			Message: "Invalid customer ID",
		})
	}

	filter := bson.M{
		"appointmentStatus": "pending",
		"paymentStatus":     "initiated",
		"customer.id":       customerObjID,
	}

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	skip := (page - 1) * limit
	findOptions := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := appointmentColl.Find(ctx, filter, findOptions, sortOptions)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(common.GetPastAppointmentsPaginationRes{
				Status:  false,
				Message: "past appointment not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(common.GetPastAppointmentsPaginationRes{
			Status:  false,
			Message: "Failed to fetch past appointment from MongoDB: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	response := common.PastAppointmentsPaginationResponse{
		Total:          0,
		PerPage:        limit,
		CurrentPage:    page,
		TotalPages:     0,
		AppointmentRes: []common.GetPastAppointmentsRes{},
	}

	for cursor.Next(ctx) {
		var appointment entity.AppointmentEntity
		err := cursor.Decode(&appointment)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(common.GetPastAppointmentsPaginationRes{
				Status:  false,
				Message: "Failed to decode past appointment data: " + err.Error(),
			})
		}

		var appointmentRes common.GetPastAppointmentsRes

		switch appointment.FacilityOrProfession {
		case "hospClinic":
			if appointment.HospitalClinic != nil {
				appointmentRes = common.GetPastAppointmentsRes{
					Id:                   appointment.Id,
					FacilityOrProfession: appointment.FacilityOrProfession,
					Image:                appointment.HospitalClinic.Information.Image,
					Name:                 appointment.HospitalClinic.Information.Name,
					PricePaid:            appointment.HospitalClinic.PricePaid,
				}
			}
		case "doctor":
			if appointment.Doctor != nil {
				appointmentRes = common.GetPastAppointmentsRes{
					Id:                   appointment.Id,
					FacilityOrProfession: appointment.FacilityOrProfession,
					Image:                appointment.Doctor.Information.Image,
					Name:                 appointment.Doctor.Information.Name,
					PricePaid:            appointment.Doctor.PricePaid,
				}
			}
		case "nurse":
			if appointment.Nurse != nil {
				appointmentRes = common.GetPastAppointmentsRes{
					Id:                   appointment.Id,
					FacilityOrProfession: appointment.FacilityOrProfession,
					Image:                appointment.Nurse.Information.Image,
					Name:                 appointment.Nurse.Information.Name,
					PricePaid:            appointment.Nurse.PricePaid,
				}
			}
		case "physiotherapist":
			if appointment.Physiotherapist != nil {
				appointmentRes = common.GetPastAppointmentsRes{
					Id:                   appointment.Id,
					FacilityOrProfession: appointment.FacilityOrProfession,
					Image:                appointment.Physiotherapist.Information.Image,
					Name:                 appointment.Physiotherapist.Information.Name,
					PricePaid:            appointment.Physiotherapist.PricePaid,
				}
			}
		case "medicalLabScientist":
			if appointment.MedicalLabScientist != nil {
				appointmentRes = common.GetPastAppointmentsRes{
					Id:                   appointment.Id,
					FacilityOrProfession: appointment.FacilityOrProfession,
					Image:                appointment.MedicalLabScientist.Information.Image,
					Name:                 appointment.MedicalLabScientist.Information.Name,
					PricePaid:            appointment.MedicalLabScientist.PricePaid,
				}
			}
		default:
			continue
		}

		if appointmentRes != (common.GetPastAppointmentsRes{}) {
			response.AppointmentRes = append(response.AppointmentRes, appointmentRes)
		}
	}

	totalCount, err := appointmentColl.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.GetPastAppointmentsPaginationRes{
			Status:  false,
			Message: "Failed to count past appointments: " + err.Error(),
		})
	}

	response.Total = int(totalCount)
	response.TotalPages = int(math.Ceil(float64(response.Total) / float64(response.PerPage)))

	finalResponse := common.GetPastAppointmentsPaginationRes{
		Status:  true,
		Message: "Sucessfully fetched data",
		Data:    response,
	}
	return c.Status(fiber.StatusOK).JSON(finalResponse)
}
