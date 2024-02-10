package fitnessCenter

import (
	"careville_backend/database"
	fitnessCenter "careville_backend/dto/admin/services/fitnessCenter"
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

// @Summary Fetch fitnessCenter With Filters
// @Description Fetch fitnessCenter With Filters
// @Tags admin fitnessCenter
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param page query int false "Page no. to fetch the products for 1"
// @Param perPage query int false "Limit of products to fetch is 15"
// @Produce json
// @Success 200 {object} fitnessCenter.GetFitnessCenterPaginationRes
// @Router /admin/healthFacility/get-fitnessCenters [get]
func FetchFitnessCenterWithPagination(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "15"))

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"role":                 "healthFacility",
		"facilityOrProfession": "fitnessCenter",
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
			return c.Status(fiber.StatusNotFound).JSON(fitnessCenter.GetFitnessCenterPaginationRes{
				Status:  false,
				Message: "fitnessCenter not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fitnessCenter.GetFitnessCenterPaginationRes{
			Status:  false,
			Message: "Failed to fetch fitnessCenter from MongoDB: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	response := fitnessCenter.FitnessCenterPaginationResponse{
		Total:            0,
		PerPage:          limit,
		CurrentPage:      page,
		TotalPages:       0,
		FitnessCenterRes: []fitnessCenter.GetFitnessCenterRes{},
	}

	for cursor.Next(ctx) {
		var service entity.ServiceEntity
		err := cursor.Decode(&service)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fitnessCenter.GetFitnessCenterPaginationRes{
				Status:  false,
				Message: "Failed to decode fitnessCenter data: " + err.Error(),
			})
		}

		fitnessCenterRes := fitnessCenter.GetFitnessCenterRes{
			Id:        service.Id,
			Email:     service.User.Email,
			FirstName: service.User.FirstName,
			LastName:  service.User.LastName,
			PhoneNumber: fitnessCenter.PhoneNumber{
				DialCode: service.User.PhoneNumber.DialCode,
				Number:   service.User.PhoneNumber.Number,
			},
			// ProfileId: service.ProfileId,
		}

		response.FitnessCenterRes = append(response.FitnessCenterRes, fitnessCenterRes)
	}

	totalCount, err := serviceColl.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fitnessCenter.GetFitnessCenterPaginationRes{
			Status:  false,
			Message: "Failed to count fitnessCenter: " + err.Error(),
		})
	}

	response.Total = int(totalCount)
	response.TotalPages = int(math.Ceil(float64(response.Total) / float64(response.PerPage)))

	finalResponse := fitnessCenter.GetFitnessCenterPaginationRes{
		Status:  true,
		Message: "Sucessfully fetched data",
		Data:    response,
	}
	return c.Status(fiber.StatusOK).JSON(finalResponse)
}
