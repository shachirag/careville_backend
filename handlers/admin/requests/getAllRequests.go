package request

import (
	"careville_backend/database"
	requests "careville_backend/dto/admin/request"
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

// @Summary Fetch requests With Filters
// @Description Fetch requests With Filters
// @Tags admin requests
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param page query int false "Page no. to fetch the products for 1"
// @Param perPage query int false "Limit of products to fetch is 15"
// @Produce json
// @Success 200 {object} requests.GetRequestsPaginationRes
// @Router /admin/requests/get-requests [get]
func FetchRequestsWithPagination(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "15"))

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"serviceStatus": "pending",
	}

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	skip := (page - 1) * limit

	projection := bson.M{
		"_id":            1,
		"user.firstName": 1,
		"user.lastName":  1,
		"user.phoneNumber": bson.M{
			"dialCode":    1,
			"number":      1,
			"countryCode": 1,
		},
		"facilityOrProfession": 1,
		"role":                 1,
		"profileId":            1,
	}

	findOptions := options.Find().SetProjection(projection).SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := serviceColl.Find(ctx, filter, findOptions, sortOptions)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(requests.GetRequestsPaginationRes{
				Status:  false,
				Message: "requests not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(requests.GetRequestsPaginationRes{
			Status:  false,
			Message: "Failed to fetch requests from MongoDB: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	response := requests.RequestsPaginationResponse{
		Total:       0,
		PerPage:     limit,
		CurrentPage: page,
		TotalPages:  0,
		Requests:    []requests.GetRequestsRes{},
	}

	for cursor.Next(ctx) {
		var service entity.ServiceEntity
		err := cursor.Decode(&service)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(requests.GetRequestsPaginationRes{
				Status:  false,
				Message: "Failed to decode requests data: " + err.Error(),
			})
		}

		requestRes := requests.GetRequestsRes{
			Id:        service.Id,
			ProfileId: service.ProfileId,
			FirstName: service.User.FirstName,
			LastName:  service.User.LastName,
			PhoneNumber: requests.PhoneNumber{
				DialCode:    service.User.PhoneNumber.DialCode,
				Number:      service.User.PhoneNumber.Number,
				CountryCode: service.User.PhoneNumber.CountryCode,
			},
			FacilityOrProfession: service.FacilityOrProfession,
			Role:                 service.Role,
		}

		response.Requests = append(response.Requests, requestRes)
	}

	totalCount, err := serviceColl.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(requests.GetRequestsPaginationRes{
			Status:  false,
			Message: "Failed to count requests: " + err.Error(),
		})
	}

	response.Total = int(totalCount)
	response.TotalPages = int(math.Ceil(float64(response.Total) / float64(response.PerPage)))

	finalResponse := requests.GetRequestsPaginationRes{
		Status:  true,
		Message: "Sucessfully fetched data",
		Data:    response,
	}
	return c.Status(fiber.StatusOK).JSON(finalResponse)
}
