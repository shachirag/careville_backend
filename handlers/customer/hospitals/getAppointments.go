package hospitals

import (
	"careville_backend/database"
	hospitals "careville_backend/dto/customer/hospitals"
	"careville_backend/entity"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Fetch appointments
// @Description Fetch appointments
// @Tags customer appointments
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param customerId query string true "customer ID"
// @Param page query int false "Page no. to fetch the products for 1"
// @Param perPage query int false "Limit of products to fetch is 15"
// @Produce json
// @Success 200 {object} hospitals.GetHospitalAppointmentsPaginationRes
// @Router /customer/healthFacility/appointment/hospital-appointments [get]
func FetchHospitalAppointmentsWithPagination(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "15"))

	appointmentColl := database.GetCollection("appointment")

	customerId := c.Query("customerId")
	customerObjID, err := primitive.ObjectIDFromHex(customerId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(hospitals.GetHospitalAppointmentsPaginationRes{
			Status:  false,
			Message: "Invalid customer ID",
		})
	}

	filter := bson.M{
		"role":                 "healthFacility",
		"facilityOrProfession": "hospClinic",
		"paymentStatus":        "initiated",
		"appointmentStatus":    "pending",
		"customer.id":          customerObjID,
	}

	projection := bson.M{
		"_id":                        1,
		"serviceId":                  1,
		"hospital.information.name":  1,
		"hospital.information.image": 1,
		"hospital.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
	}

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	skip := (page - 1) * limit
	findOptions := options.Find().SetProjection(projection).SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := appointmentColl.Find(ctx, filter, findOptions, sortOptions)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(hospitals.GetHospitalAppointmentsPaginationRes{
				Status:  false,
				Message: "appointment not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(hospitals.GetHospitalAppointmentsPaginationRes{
			Status:  false,
			Message: "Failed to fetch appointment from MongoDB: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	response := hospitals.HospitalAppointmentsPaginationResponse{
		Total:          0,
		PerPage:        limit,
		CurrentPage:    page,
		TotalPages:     0,
		AppointmentRes: []hospitals.GetHospitalAppointmentsRes{},
	}

	for cursor.Next(ctx) {
		var appointment entity.AppointmentEntity
		err := cursor.Decode(&appointment)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(hospitals.GetHospitalAppointmentsPaginationRes{
				Status:  false,
				Message: "Failed to decode appointment data: " + err.Error(),
			})
		}

		if appointment.HospitalClinic != nil {
			appointmentRes := hospitals.GetHospitalAppointmentsRes{
				Id:         appointment.Id,
				HospitalId: appointment.ServiceID,
				Image:      appointment.HospitalClinic.Information.Image,
				Name:       appointment.HospitalClinic.Information.Name,
				Address:    hospitals.Address(appointment.HospitalClinic.Information.Address),
			}

			response.AppointmentRes = append(response.AppointmentRes, appointmentRes)
		}
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(hospitals.GetHospitalAppointmentsPaginationRes{
			Status:  false,
			Message: "Failed to count appointments: " + err.Error(),
		})
	}

	totalCount, err := appointmentColl.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(hospitals.GetHospitalAppointmentsPaginationRes{
			Status:  false,
			Message: "Failed to count appointments: " + err.Error(),
		})
	}

	response.Total = int(totalCount)
	response.TotalPages = int(math.Ceil(float64(response.Total) / float64(response.PerPage)))

	finalResponse := hospitals.GetHospitalAppointmentsPaginationRes{
		Status:  true,
		Message: "Sucessfully fetched data",
		Data:    response,
	}
	return c.Status(fiber.StatusOK).JSON(finalResponse)
}
