package pharmacy

import (
	"careville_backend/database"
	pharmacy "careville_backend/dto/admin/services/pharmacy"
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

// @Summary Fetch pharmacy With Filters
// @Description Fetch pharmacy With Filters
// @Tags admin pharmacy
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param page query int false "Page no. to fetch the products for 1"
// @Param perPage query int false "Limit of products to fetch is 15"
// @Produce json
// @Success 200 {object} pharmacy.GetPharmacyPaginationRes
// @Router /admin/healthFacility/get-pharmacies [get]
func FetchPharmacyWithPagination(c *fiber.Ctx) error {

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "15"))

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"role":                 "healthFacility",
		"facilityOrProfession": "pharmacy",
		"serviceStatus":        "approved",
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
		"profileId": 1,
	}

	findOptions := options.Find().SetProjection(projection).SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := serviceColl.Find(ctx, filter, findOptions, sortOptions)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(pharmacy.GetPharmacyPaginationRes{
				Status:  false,
				Message: "pharmacy not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(pharmacy.GetPharmacyPaginationRes{
			Status:  false,
			Message: "Failed to fetch pharmacy from MongoDB: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	response := pharmacy.PharmacyPaginationResponse{
		Total:       0,
		PerPage:     limit,
		CurrentPage: page,
		TotalPages:  0,
		PharmacyRes: []pharmacy.GetPharmacyRes{},
	}

	for cursor.Next(ctx) {
		var service entity.ServiceEntity
		err := cursor.Decode(&service)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(pharmacy.GetPharmacyPaginationRes{
				Status:  false,
				Message: "Failed to decode pharmacy data: " + err.Error(),
			})
		}

		pharmacyRes := pharmacy.GetPharmacyRes{
			Id:        service.Id,
			Email:     service.User.Email,
			FirstName: service.User.FirstName,
			LastName:  service.User.LastName,
			PhoneNumber: pharmacy.PhoneNumber{
				DialCode: service.User.PhoneNumber.DialCode,
				Number:   service.User.PhoneNumber.Number,
			},
			ProfileId: service.ProfileId,
		}

		response.PharmacyRes = append(response.PharmacyRes, pharmacyRes)
	}

	totalCount, err := serviceColl.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pharmacy.GetPharmacyPaginationRes{
			Status:  false,
			Message: "Failed to count pharmacy: " + err.Error(),
		})
	}

	response.Total = int(totalCount)
	response.TotalPages = int(math.Ceil(float64(response.Total) / float64(response.PerPage)))

	finalResponse := pharmacy.GetPharmacyPaginationRes{
		Status:  true,
		Message: "Sucessfully fetched data",
		Data:    response,
	}
	return c.Status(fiber.StatusOK).JSON(finalResponse)
}
