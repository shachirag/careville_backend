package laboratory

import (
	"careville_backend/database"
	laboratory "careville_backend/dto/customer/laboratories"
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
// @Success 200 {object} laboratory.GetLaboratoryAppointmentsPaginationRes
// @Router /customer/healthFacility/appointment/laboratory-appointments [get]
func FetchLaboratoryAppointmentsWithPagination(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "15"))

	appointmentColl := database.GetCollection("appointment")

	customerId := c.Query("customerId")
	customerObjID, err := primitive.ObjectIDFromHex(customerId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(laboratory.GetLaboratoryAppointmentsPaginationRes{
			Status:  false,
			Message: "Invalid customer ID",
		})
	}

	filter := bson.M{
		"role":                 "healthFacility",
		"facilityOrProfession": "laboratory",
		"appointmentStatus":    "pending",
		"customer.id":          customerObjID,
	}

	projection := bson.M{
		"_id":                          1,
		"serviceId":                    1,
		"laboratory.information.name":  1,
		"laboratory.information.image": 1,
		"laboratory.information.address": bson.M{
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
			return c.Status(fiber.StatusNotFound).JSON(laboratory.GetLaboratoryAppointmentsPaginationRes{
				Status:  false,
				Message: "appointment not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(laboratory.GetLaboratoryAppointmentsPaginationRes{
			Status:  false,
			Message: "Failed to fetch appointment from MongoDB: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	response := laboratory.LaboratoryAppointmentsPaginationResponse{
		Total:          0,
		PerPage:        limit,
		CurrentPage:    page,
		TotalPages:     0,
		AppointmentRes: []laboratory.GetLaboratoryAppointmentsRes{},
	}

	for cursor.Next(ctx) {
		var appointment entity.AppointmentEntity
		err := cursor.Decode(&appointment)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(laboratory.GetLaboratoryAppointmentsPaginationRes{
				Status:  false,
				Message: "Failed to decode appointment data: " + err.Error(),
			})
		}

		if appointment.Laboratory != nil {
			appointmentRes := laboratory.GetLaboratoryAppointmentsRes{
				Id:           appointment.Id,
				LaboratoryID: appointment.ServiceID,
				Image:        appointment.Laboratory.Information.Image,
				Name:         appointment.Laboratory.Investigation.Name,
				Address:      laboratory.Address(appointment.Laboratory.Information.Address),
			}

			response.AppointmentRes = append(response.AppointmentRes, appointmentRes)
		}
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(laboratory.GetLaboratoryAppointmentsPaginationRes{
			Status:  false,
			Message: "Failed to count appointments: " + err.Error(),
		})
	}

	totalCount, err := appointmentColl.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(laboratory.GetLaboratoryAppointmentsPaginationRes{
			Status:  false,
			Message: "Failed to count appointments: " + err.Error(),
		})
	}

	response.Total = int(totalCount)
	response.TotalPages = int(math.Ceil(float64(response.Total) / float64(response.PerPage)))

	finalResponse := laboratory.GetLaboratoryAppointmentsPaginationRes{
		Status:  true,
		Message: "Sucessfully fetched data",
		Data:    response,
	}
	return c.Status(fiber.StatusOK).JSON(finalResponse)
}
