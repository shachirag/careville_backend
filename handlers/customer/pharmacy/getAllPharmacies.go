package pharmacy

import (
	"careville_backend/database"
	pharmacy "careville_backend/dto/customer/pharmacy"
	"careville_backend/entity"
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

// @Summary Get pharmacy
// @Tags customer pharmacy
// @Description Get pharmacy
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param search query string false "Filter pharmacy by search"
// @Param long query float64 false "Longitude for memories sorting (required for distance sorting)"
// @Param lat query float64 false "Latitude for memories sorting (required for distance sorting)"
// @Produce json
// @Success 200 {object} pharmacy.GetPharmacyResponseDto
// @Router /customer/healthFacility/get-pharmacies [get]
func GetPharmacy(c *fiber.Ctx) error {

	var (
		serviceColl = database.GetCollection("service")
	)

	searchTitle := c.Query("search", "")

	var lat, long float64
	latParam := c.Query("lat")
	longParam := c.Query("long")
	var err error

	if latParam != "" && longParam != "" {
		lat, err = strconv.ParseFloat(latParam, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(pharmacy.GetPharmacyResponseDto{
				Status:  false,
				Message: "Invalid latitude format",
			})
		}

		long, err = strconv.ParseFloat(longParam, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(pharmacy.GetPharmacyResponseDto{
				Status:  false,
				Message: "Invalid longitude format",
			})
		}
	}

	filter := bson.M{
		"role":                 "healthFacility",
		"facilityOrProfession": "pharmacy",
	}

	maxDistance := 50000

	if latParam != "" && longParam != "" {
		filter["pharmacy.information.address"] = bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{long, lat},
				},
				"$maxDistance": maxDistance,
			},
		}
	}

	if searchTitle != "" {
		filter["pharmacy.information.name"] = bson.M{"$regex": searchTitle, "$options": "i"}
	}

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	projection := bson.M{
		"pharmacy.information.name":  1,
		"pharmacy.information.image": 1,
		"_id":                        1,
		"pharmacy.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
		"pharmacy.review.avgRating": 1,
	}

	findOptions := options.Find().SetProjection(projection)

	cursor, err := serviceColl.Find(ctx, filter, sortOptions, findOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pharmacy.GetPharmacyResponseDto{
			Status:  false,
			Message: "Failed to fetch pharmacy data: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var pharmacy1Data []pharmacy.GetPharmacyRes
	for cursor.Next(ctx) {
		var pharmacy1 entity.ServiceEntity
		if err := cursor.Decode(&pharmacy1); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(pharmacy.GetPharmacyResponseDto{
				Status:  false,
				Message: "Failed to decode pharmacy data: " + err.Error(),
			})
		}
		if pharmacy1.Pharmacy != nil {
			pharmacy1Data = append(pharmacy1Data, pharmacy.GetPharmacyRes{
				Id:    pharmacy1.Id,
				Image: pharmacy1.Pharmacy.Information.Image,
				Name:  pharmacy1.Pharmacy.Information.Name,
				Address: pharmacy.Address{
					Coordinates: pharmacy1.Pharmacy.Information.Address.Coordinates,
					Type:        pharmacy1.Pharmacy.Information.Address.Type,
					Add:         pharmacy1.Pharmacy.Information.Address.Add,
				},
				AvgRating: pharmacy1.Pharmacy.Review.AvgRating,
			})
		}
	}

	if len(pharmacy1Data) == 0 {
		return c.Status(fiber.StatusOK).JSON(pharmacy.GetPharmacyResponseDto{
			Status:  false,
			Message: "No Pharmacy data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(pharmacy.GetPharmacyResponseDto{
		Status:  true,
		Message: "Successfully fetched pharmacy data.",
		Data:    pharmacy1Data,
	})
}
