package physiotherapist

import (
	"careville_backend/database"
	physiotherapist "careville_backend/dto/admin/services/physiotherapist"
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

// @Summary Fetch physiotherapist With Filters
// @Description Fetch physiotherapist With Filters
// @Tags admin physiotherapist
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param page query int false "Page no. to fetch the products for 1"
// @Param perPage query int false "Limit of products to fetch is 15"
// @Produce json
// @Success 200 {object} physiotherapist.GetPhysiotherapistPaginationRes
// @Router /admin/healthProfessional/get-physiotherapists [get]
func FetchPhysiotherapistWithPagination(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "15"))

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"role":                 "healthProfessional",
		"facilityOrProfession": "physiotherapist",
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
			return c.Status(fiber.StatusNotFound).JSON(physiotherapist.GetPhysiotherapistPaginationRes{
				Status:  false,
				Message: "physiotherapist not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(physiotherapist.GetPhysiotherapistPaginationRes{
			Status:  false,
			Message: "Failed to fetch physiotherapist from MongoDB: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	response := physiotherapist.PhysiotherapistPaginationResponse{
		Total:              0,
		PerPage:            limit,
		CurrentPage:        page,
		TotalPages:         0,
		PhysiotherapistRes: []physiotherapist.GetPhysiotherapistRes{},
	}

	for cursor.Next(ctx) {
		var service entity.ServiceEntity
		err := cursor.Decode(&service)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(physiotherapist.GetPhysiotherapistPaginationRes{
				Status:  false,
				Message: "Failed to decode physiotherapist data: " + err.Error(),
			})
		}

		nurseRes := physiotherapist.GetPhysiotherapistRes{
			Id:        service.Id,
			Email:     service.User.Email,
			FirstName: service.User.FirstName,
			LastName:  service.User.LastName,
			PhoneNumber: physiotherapist.PhoneNumber{
				DialCode: service.User.PhoneNumber.DialCode,
				Number:   service.User.PhoneNumber.Number,
			},
			// ProfileId: service.ProfileId,
		}

		response.PhysiotherapistRes = append(response.PhysiotherapistRes, nurseRes)
	}

	totalCount, err := serviceColl.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(physiotherapist.GetPhysiotherapistPaginationRes{
			Status:  false,
			Message: "Failed to count physiotherapist: " + err.Error(),
		})
	}

	response.Total = int(totalCount)
	response.TotalPages = int(math.Ceil(float64(response.Total) / float64(response.PerPage)))

	finalResponse := physiotherapist.GetPhysiotherapistPaginationRes{
		Status:  true,
		Message: "Sucessfully fetched data",
		Data:    response,
	}
	return c.Status(fiber.StatusOK).JSON(finalResponse)
}
