package medicalLabScientist

import (
	"careville_backend/database"
	medicalLabScientist "careville_backend/dto/customer/medicalLabScientist"
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

// @Summary Fetch medicalLabScientists With Filters
// @Description Fetch medicalLabScientists With Filters
// @Tags customer medicalLabScientist
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param page query int false "Page no. to fetch the products for 1"
// @Param perPage query int false "Limit of products to fetch is 15"
// @Param long query float64 false "Longitude for memories sorting (required for distance sorting)"
// @Param lat query float64 false "Latitude for memories sorting (required for distance sorting)"
// @Param search query string false "Filter laboratory by search"
// @Produce json
// @Success 200 {object} medicalLabScientist.GetMedicalLabScientistPaginationRes
// @Router /customer/healthProfessional/get-medicalLabScientists [get]
func FetchMedicalLabScientistsWithPagination(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "15"))

	var lat, long float64
	latParam := c.Query("lat")
	longParam := c.Query("long")
	var err error

	if latParam != "" && longParam != "" {
		lat, err = strconv.ParseFloat(latParam, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(medicalLabScientist.GetMedicalLabScientistPaginationRes{
				Status:  false,
				Message: "Invalid latitude format",
			})
		}

		long, err = strconv.ParseFloat(longParam, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(medicalLabScientist.GetMedicalLabScientistPaginationRes{
				Status:  false,
				Message: "Invalid longitude format",
			})
		}
	}

	searchTitle := c.Query("search", "")

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"role":                 "healthProfessional",
		"facilityOrProfession": "medicalLabScientist",
	}

	if latParam != "" && longParam != "" {
		filter["medicalLabScientist.information.address"] = bson.M{
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
		filter["medicalLabScientist.information.name"] = bson.M{"$regex": searchTitle, "$options": "i"}
	}

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	skip := (page - 1) * limit

	projection := bson.M{
		"medicalLabScientist.information.name":  1,
		"medicalLabScientist.information.image": 1,
		"medicalLabScientist.information.id":    1,
		"avgRating":                             1,
	}

	findOptions := options.Find().SetProjection(projection).SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := serviceColl.Find(ctx, filter, findOptions, sortOptions)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(medicalLabScientist.GetMedicalLabScientistPaginationRes{
				Status:  false,
				Message: "medicalLabScientist not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.GetMedicalLabScientistPaginationRes{
			Status:  false,
			Message: "Failed to fetch medicalLabScientist from MongoDB: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	response := medicalLabScientist.MedicalLabScientistPaginationResponse{
		Total:                  0,
		PerPage:                limit,
		CurrentPage:            page,
		TotalPages:             0,
		MedicalLabScientistRes: []medicalLabScientist.GetMedicalLabScientistRes{},
	}

	for cursor.Next(ctx) {
		var service entity.ServiceEntity
		err := cursor.Decode(&service)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.GetMedicalLabScientistPaginationRes{
				Status:  false,
				Message: "Failed to decode medicalLabScientist data: " + err.Error(),
			})
		}

		// Check if hospClinic is not nil before accessing its properties
		if service.MedicalLabScientist != nil {
			medicalLabScientistRes := medicalLabScientist.GetMedicalLabScientistRes{
				Id:        service.Id,
				Image:     service.MedicalLabScientist.Information.Image,
				Name:      service.MedicalLabScientist.Information.Name,
				AvgRating: service.AvgRating,
			}

			response.MedicalLabScientistRes = append(response.MedicalLabScientistRes, medicalLabScientistRes)
		}
	}

	totalCount, err := serviceColl.CountDocuments(ctx, bson.M{
		"role":                                 "healthProfessional",
		"facilityOrProfession":                 "medicalLabScientist",
		"medicalLabScientist.information.name": bson.M{"$regex": searchTitle, "$options": "i"},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.GetMedicalLabScientistPaginationRes{
			Status:  false,
			Message: "Failed to count medicalLabScientist: " + err.Error(),
		})
	}

	response.Total = int(totalCount)
	response.TotalPages = int(math.Ceil(float64(response.Total) / float64(response.PerPage)))

	finalResponse := medicalLabScientist.GetMedicalLabScientistPaginationRes{
		Status:  true,
		Message: "Sucessfully fetched data",
		Data:    response,
	}
	return c.Status(fiber.StatusOK).JSON(finalResponse)
}
