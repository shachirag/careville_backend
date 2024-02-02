package hospitals

import (
	"careville_backend/database"
	hospitals "careville_backend/dto/customer/hospitals"
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

// @Summary Fetch hospitals With Filters
// @Description Fetch hospitals With Filters
// @Tags customer hospitals
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param page query int false "Page no. to fetch the products for 1"
// @Param perPage query int false "Limit of products to fetch is 15"
// @Param long query float64 false "Longitude for memories sorting (required for distance sorting)"
// @Param lat query float64 false "Latitude for memories sorting (required for distance sorting)"
// @Param search query string false "Filter hospitals by search"
// @Produce json
// @Success 200 {object} hospitals.GetHospitalsPaginationRes
// @Router /customer/healthFacility/get-hospitals [get]
func FetchHospitalsWithPagination(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "15"))

	var lat, long float64
	latParam := c.Query("lat")
	longParam := c.Query("long")
	var err error

	if latParam != "" && longParam != "" {
		lat, err = strconv.ParseFloat(latParam, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(hospitals.GetHospitalsPaginationRes{
				Status:  false,
				Message: "Invalid latitude format",
			})
		}

		long, err = strconv.ParseFloat(longParam, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(hospitals.GetHospitalsPaginationRes{
				Status:  false,
				Message: "Invalid longitude format",
			})
		}
	}

	searchTitle := c.Query("search", "")

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"role":                 "healthFacility",
		"facilityOrProfession": "hospClinic",
		"hospClinic.information.address": bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{long, lat},
				},
				"$maxDistance": 20000,
			},
		},
	}

	if searchTitle != "" {
		filter["hospClinic.information.name"] = bson.M{"$regex": searchTitle, "$options": "i"}
	}

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	skip := (page - 1) * limit

	projection := bson.M{
		"hospClinic.information.name":  1,
		"hospClinic.information.image": 1,
		"hospClinic.information.id":    1,
		"hospClinic.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
	}

	findOptions := options.Find().SetProjection(projection).SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := serviceColl.Find(ctx, filter, findOptions, sortOptions)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(hospitals.GetHospitalsPaginationRes{
				Status:  false,
				Message: "hospital not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(hospitals.GetHospitalsPaginationRes{
			Status:  false,
			Message: "Failed to fetch hospital from MongoDB: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	response := hospitals.HospitalsPaginationResponse{
		Total:       0,
		PerPage:     limit,
		CurrentPage: page,
		TotalPages:  0,
		HospitalRes: []hospitals.GetHospitalsRes{},
	}

	for cursor.Next(ctx) {
		var service entity.ServiceEntity
		err := cursor.Decode(&service)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(hospitals.GetHospitalsPaginationRes{
				Status:  false,
				Message: "Failed to decode hospital data: " + err.Error(),
			})
		}

		// Check if hospClinic is not nil before accessing its properties
		if service.HospClinic != nil {
			hospitalRes := hospitals.GetHospitalsRes{
				Id:    service.Id,
				Image: service.HospClinic.Information.Image,
				Name:  service.HospClinic.Information.Name,
				Address: hospitals.Address{
					Coordinates: service.HospClinic.Information.Address.Coordinates,
					Type:        service.HospClinic.Information.Address.Type,
					Add:         service.HospClinic.Information.Address.Add,
				},
			}

			response.HospitalRes = append(response.HospitalRes, hospitalRes)
		}
	}

	totalCount, err := serviceColl.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(hospitals.GetHospitalsPaginationRes{
			Status:  false,
			Message: "Failed to count hospitals: " + err.Error(),
		})
	}

	response.Total = int(totalCount)
	response.TotalPages = int(math.Ceil(float64(response.Total) / float64(response.PerPage)))

	finalResponse := hospitals.GetHospitalsPaginationRes{
		Status:  true,
		Message: "Sucessfully fetched data",
		Data:    response,
	}
	return c.Status(fiber.StatusOK).JSON(finalResponse)
}
