package laboratory

import (
	"careville_backend/database"
	laboratory "careville_backend/dto/admin/services/laboratories"
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

// @Summary Fetch laboratory With Filters
// @Description Fetch laboratory With Filters
// @Tags admin laboratory
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param page query int false "Page no. to fetch the products for 1"
// @Param perPage query int false "Limit of products to fetch is 15"
// @Produce json
// @Success 200 {object} laboratory.GetLaboratoryPaginationRes
// @Router /admin/healthFacility/get-laboratories [get]
func FetchLaboratoriesWithPagination(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "15"))

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"role":                 "healthFacility",
		"facilityOrProfession": "laboratory",
	}

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	skip := (page - 1) * limit

	projection := bson.M{
		"_id":            1,
		"user.firstName": 1,
		"user.lastName":  1,
		"user.email":     1,
		"user.phoneNumber": bson.M{
			"dialCode": 1,
			"number":   1,
		},
	}

	findOptions := options.Find().SetProjection(projection).SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := serviceColl.Find(ctx, filter, findOptions, sortOptions)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(laboratory.GetLaboratoryPaginationRes{
				Status:  false,
				Message: "laboratories not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(laboratory.GetLaboratoryPaginationRes{
			Status:  false,
			Message: "Failed to fetch laboratories from MongoDB: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	response := laboratory.LaboratoryPaginationResponse{
		Total:         0,
		PerPage:       limit,
		CurrentPage:   page,
		TotalPages:    0,
		LaboratoryRes: []laboratory.GetLaboratoryRes{},
	}

	for cursor.Next(ctx) {
		var service entity.ServiceEntity
		err := cursor.Decode(&service)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(laboratory.GetLaboratoryPaginationRes{
				Status:  false,
				Message: "Failed to decode laboratories data: " + err.Error(),
			})
		}

		laboratoryRes := laboratory.GetLaboratoryRes{
			Id:        service.Id,
			Email:     service.User.Email,
			FirstName: service.User.FirstName,
			LastName:  service.User.LastName,
			PhoneNumber: laboratory.PhoneNumber{
				DialCode: service.User.PhoneNumber.DialCode,
				Number:   service.User.PhoneNumber.Number,
			},
		}

		response.LaboratoryRes = append(response.LaboratoryRes, laboratoryRes)
	}

	totalCount, err := serviceColl.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(laboratory.GetLaboratoryPaginationRes{
			Status:  false,
			Message: "Failed to count laboratories: " + err.Error(),
		})
	}

	response.Total = int(totalCount)
	response.TotalPages = int(math.Ceil(float64(response.Total) / float64(response.PerPage)))

	finalResponse := laboratory.GetLaboratoryPaginationRes{
		Status:  true,
		Message: "Sucessfully fetched data",
		Data:    response,
	}
	return c.Status(fiber.StatusOK).JSON(finalResponse)
}
