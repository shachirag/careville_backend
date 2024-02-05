package nurse

import (
	"careville_backend/database"
	nurse "careville_backend/dto/customer/nurse"
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

// @Summary Fetch nurse With Filters
// @Description Fetch nurse With Filters
// @Tags customer nurse
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param page query int false "Page no. to fetch the products for 1"
// @Param perPage query int false "Limit of products to fetch is 15"
// @Param long query float64 false "Longitude for memories sorting (required for distance sorting)"
// @Param lat query float64 false "Latitude for memories sorting (required for distance sorting)"
// @Param search query string false "Filter nurse by search"
// @Produce json
// @Success 200 {object} nurse.GetNursePaginationRes
// @Router /customer/healthProfessional/get-nurses [get]
func FetchNurseWithPagination(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "15"))

	var lat, long float64
	latParam := c.Query("lat")
	longParam := c.Query("long")
	var err error

	if latParam != "" && longParam != "" {
		lat, err = strconv.ParseFloat(latParam, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(nurse.GetNursePaginationRes{
				Status:  false,
				Message: "Invalid latitude format",
			})
		}

		long, err = strconv.ParseFloat(longParam, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(nurse.GetNursePaginationRes{
				Status:  false,
				Message: "Invalid longitude format",
			})
		}
	}

	searchTitle := c.Query("search", "")

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"role":                 "healthProfessional",
		"facilityOrProfession": "nurse",
	}

	if latParam != "" && longParam != "" {
		filter["nurse.information.address"] = bson.M{
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
		filter["nurse.information.name"] = bson.M{"$regex": searchTitle, "$options": "i"}
	}

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	skip := (page - 1) * limit

	projection := bson.M{
		"nurse.information.name":  1,
		"nurse.information.image": 1,
		"nurse.information.id":    1,
	}

	findOptions := options.Find().SetProjection(projection).SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := serviceColl.Find(ctx, filter, findOptions, sortOptions)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(nurse.GetNursePaginationRes{
				Status:  false,
				Message: "nurse not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(nurse.GetNursePaginationRes{
			Status:  false,
			Message: "Failed to fetch nurse from MongoDB: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	response := nurse.NursePaginationResponse{
		Total:       0,
		PerPage:     limit,
		CurrentPage: page,
		TotalPages:  0,
		NurseRes:    []nurse.GetNurseRes{},
	}

	for cursor.Next(ctx) {
		var service entity.ServiceEntity
		err := cursor.Decode(&service)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(nurse.GetNursePaginationRes{
				Status:  false,
				Message: "Failed to decode nurse data: " + err.Error(),
			})
		}

		// Check if hospClinic is not nil before accessing its properties
		if service.Nurse != nil {
			nurseRes := nurse.GetNurseRes{
				Id:    service.Id,
				Image: service.Nurse.Information.Image,
				Name:  service.Nurse.Information.Name,
			}

			response.NurseRes = append(response.NurseRes, nurseRes)
		}
	}

	totalCount, err := serviceColl.CountDocuments(ctx, bson.M{
		"role":                           "healthProfessional",
		"facilityOrProfession":           "nurse",
		"nurse.information.name": bson.M{"$regex": searchTitle, "$options": "i"},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(nurse.GetNursePaginationRes{
			Status:  false,
			Message: "Failed to count nurse: " + err.Error(),
		})
	}

	response.Total = int(totalCount)
	response.TotalPages = int(math.Ceil(float64(response.Total) / float64(response.PerPage)))

	finalResponse := nurse.GetNursePaginationRes{
		Status:  true,
		Message: "Sucessfully fetched data",
		Data:    response,
	}
	return c.Status(fiber.StatusOK).JSON(finalResponse)
}
