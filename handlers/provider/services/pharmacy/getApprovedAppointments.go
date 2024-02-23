package pharmacy

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	pharmacy "careville_backend/dto/provider/services"
	"careville_backend/entity"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Fetch appointments
// @Description Fetch drugs
// @Tags provider appointments
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param page query int false "Page no. to fetch the products for 1"
// @Param perPage query int false "Limit of products to fetch is 15"
// @Produce json
// @Success 200 {object} pharmacy.GetPharmacyAppointmentsPaginationRes
// @Router /provider/services/appointment/pharmacy-drugs [get]
func FetchPharmacyApprovedDrugsWithPagination(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "15"))

	appointmentColl := database.GetCollection("appointment")

	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	filter := bson.M{
		"role":                 "healthFacility",
		"facilityOrProfession": "pharmacy",
		"paymentStatus":        "success",
		"appointmentStatus":    "approved",
		"serviceId":            providerData.ProviderId,
	}

	projection := bson.M{
		"_id":                  1,
		"customer.id":          1,
		"customer.firstName":   1,
		"customer.lastName":    1,
		"facilityOrProfession": 1,
		"createdAt":            1,
	}

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	skip := (page - 1) * limit
	findOptions := options.Find().SetProjection(projection).SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := appointmentColl.Find(ctx, filter, findOptions, sortOptions)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(pharmacy.GetPharmacyAppointmentsPaginationRes{
				Status:  false,
				Message: "appointment not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(pharmacy.GetPharmacyAppointmentsPaginationRes{
			Status:  false,
			Message: "Failed to fetch appointment from MongoDB: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	response := pharmacy.PharmacyAppointmentsPaginationResponse{
		Total:       0,
		PerPage:     limit,
		CurrentPage: page,
		TotalPages:  0,
		Drugs:       []pharmacy.GetPharmacyAppointmentsRes{},
	}

	for cursor.Next(ctx) {
		var appointment entity.AppointmentEntity
		err := cursor.Decode(&appointment)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(pharmacy.GetPharmacyAppointmentsPaginationRes{
				Status:  false,
				Message: "Failed to decode appointment data: " + err.Error(),
			})
		}

		appointmentRes := pharmacy.GetPharmacyAppointmentsRes{
			Id:                   appointment.Id,
			CustomrId:            appointment.Customer.ID,
			FacilityOrProfession: appointment.FacilityOrProfession,
			FirstName:            appointment.Customer.FirstName,
			LastName:             appointment.Customer.LastName,
			CreatedAt:            appointment.CreatedAt,
		}

		response.Drugs = append(response.Drugs, appointmentRes)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pharmacy.GetPharmacyAppointmentsPaginationRes{
			Status:  false,
			Message: "Failed to count appointments: " + err.Error(),
		})
	}

	totalCount, err := appointmentColl.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pharmacy.GetPharmacyAppointmentsPaginationRes{
			Status:  false,
			Message: "Failed to count appointments: " + err.Error(),
		})
	}

	response.Total = int(totalCount)
	response.TotalPages = int(math.Ceil(float64(response.Total) / float64(response.PerPage)))

	finalResponse := pharmacy.GetPharmacyAppointmentsPaginationRes{
		Status:  true,
		Message: "Sucessfully fetched data",
		Data:    response,
	}
	return c.Status(fiber.StatusOK).JSON(finalResponse)
}
