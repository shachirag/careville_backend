package doctorProfession

import (
	"careville_backend/database"
	doctorProfession "careville_backend/dto/customer/doctorProfession"
	"careville_backend/entity"
	"context"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

// @Summary Fetch doctorProfession With Filters
// @Description Fetch doctorProfession With Filters
// @Tags customer doctorProfession
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param page query int false "Page no. to fetch the products for 1"
// @Param perPage query int false "Limit of products to fetch is 15"
// @Param long query float64 false "Longitude for memories sorting (required for distance sorting)"
// @Param lat query float64 false "Latitude for memories sorting (required for distance sorting)"
// @Param search query string false "Filter doctorProfession by search"
// @Produce json
// @Success 200 {object} doctorProfession.GetDoctorProfessionPaginationRes
// @Router /customer/healthProfessional/get-doctors [get]
func FetchDoctorsWithPagination(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "15"))

	var lat, long float64
	latParam := c.Query("lat")
	longParam := c.Query("long")
	var err error

	if latParam != "" && longParam != "" {
		lat, err = strconv.ParseFloat(latParam, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(doctorProfession.GetDoctorProfessionPaginationRes{
				Status:  false,
				Message: "Invalid latitude format",
			})
		}

		long, err = strconv.ParseFloat(longParam, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(doctorProfession.GetDoctorProfessionPaginationRes{
				Status:  false,
				Message: "Invalid longitude format",
			})
		}
	}

	searchTitle := c.Query("search", "")

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"role":                 "healthProfessional",
		"facilityOrProfession": "doctor",
	}

	if latParam != "" && longParam != "" {
		filter["doctor.information.address"] = bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{long, lat},
				},
				"$maxDistance": 20000,
			},
		}
	}

	if searchTitle != "" {
		filter["doctor.information.name"] = bson.M{"$regex": searchTitle, "$options": "i"}
	}

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	skip := (page - 1) * limit

	projection := bson.M{
		"doctor.information.name":              1,
		"doctor.information.image":             1,
		"doctor.information.id":                1,
		"doctor.additionalServices.speciality": 1,
	}

	findOptions := options.Find().SetProjection(projection).SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := serviceColl.Find(ctx, filter, findOptions, sortOptions)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(doctorProfession.GetDoctorProfessionPaginationRes{
				Status:  false,
				Message: "doctor not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(doctorProfession.GetDoctorProfessionPaginationRes{
			Status:  false,
			Message: "Failed to fetch doctor from MongoDB: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	response := doctorProfession.DoctorProfessionPaginationResponse{
		Total:               0,
		PerPage:             limit,
		CurrentPage:         page,
		TotalPages:          0,
		DoctorProfessionRes: []doctorProfession.GetDoctorProfessionRes{},
	}

	for cursor.Next(ctx) {
		var service entity.ServiceEntity
		err := cursor.Decode(&service)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(doctorProfession.GetDoctorProfessionPaginationRes{
				Status:  false,
				Message: "Failed to decode doctor data: " + err.Error(),
			})
		}

		// Check if hospClinic is not nil before accessing its properties
		if service.Doctor != nil {
			medicalLabScientistRes := doctorProfession.GetDoctorProfessionRes{
				Id:         service.Id,
				Image:      service.Doctor.Information.Image,
				Name:       service.Doctor.Information.Name,
				Speciality: service.Doctor.AdditionalServices.Speciality,
			}

			response.DoctorProfessionRes = append(response.DoctorProfessionRes, medicalLabScientistRes)
		}
	}

	totalCount, err := serviceColl.CountDocuments(ctx, bson.M{
		"role":                        "healthProfessional",
		"facilityOrProfession":        "doctor",
		"doctor.information.name": bson.M{"$regex": searchTitle, "$options": "i"},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(doctorProfession.GetDoctorProfessionPaginationRes{
			Status:  false,
			Message: "Failed to count doctor: " + err.Error(),
		})
	}

	response.Total = int(totalCount)
	response.TotalPages = int(math.Ceil(float64(response.Total) / float64(response.PerPage)))

	finalResponse := doctorProfession.GetDoctorProfessionPaginationRes{
		Status:  true,
		Message: "Sucessfully fetched data",
		Data:    response,
	}
	return c.Status(fiber.StatusOK).JSON(finalResponse)
}
