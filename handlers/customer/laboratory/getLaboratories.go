package laboratory

import (
	"careville_backend/database"
	laboratory "careville_backend/dto/customer/laboratories"
	"careville_backend/entity"
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

// @Summary Get laboratory
// @Tags customer laboratory
// @Description Get laboratory
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param search query string false "Filter laboratory by search"
// @Param long query float64 false "Longitude for memories sorting (required for distance sorting)"
// @Param lat query float64 false "Latitude for memories sorting (required for distance sorting)"
// @Produce json
// @Success 200 {object} laboratory.GetLaboratoryResponseDto
// @Router /customer/healthFacility/get-laboratories [get]
func Getlaboratory(c *fiber.Ctx) error {

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
			return c.Status(fiber.StatusBadRequest).JSON(laboratory.GetLaboratoryResponseDto{
				Status:  false,
				Message: "Invalid latitude format",
			})
		}

		long, err = strconv.ParseFloat(longParam, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(laboratory.GetLaboratoryResponseDto{
				Status:  false,
				Message: "Invalid longitude format",
			})
		}
	}

	filter := bson.M{
		"role":                 "healthFacility",
		"facilityOrProfession": "laboratory",
	}

	maxDistance := 50000

	if latParam != "" && longParam != "" {
		filter["laboratory.information.address"] = bson.M{
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
		filter["laboratory.information.name"] = bson.M{"$regex": searchTitle, "$options": "i"}
	}

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	projection := bson.M{
		"laboratory.information.name":  1,
		"laboratory.information.image": 1,
		"_id":                          1,
		"laboratory.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
		"laboratory.review.avgRating": 1,
	}
	findOptions := options.Find().SetProjection(projection)

	cursor, err := serviceColl.Find(ctx, filter, sortOptions, findOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(laboratory.GetLaboratoryResponseDto{
			Status:  false,
			Message: "Failed to fetch laboratory data: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var laboratoryData []laboratory.GetLaboratoryRes
	for cursor.Next(ctx) {
		var laboratory1 entity.ServiceEntity
		if err := cursor.Decode(&laboratory1); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(laboratory.GetLaboratoryResponseDto{
				Status:  false,
				Message: "Failed to decode laboratory data: " + err.Error(),
			})
		}
		if laboratory1.Laboratory != nil {
			laboratoryData = append(laboratoryData, laboratory.GetLaboratoryRes{
				Id:    laboratory1.Id,
				Image: laboratory1.Laboratory.Information.Image,
				Name:  laboratory1.Laboratory.Information.Name,
				Address: laboratory.Address{
					Coordinates: laboratory1.Laboratory.Information.Address.Coordinates,
					Type:        laboratory1.Laboratory.Information.Address.Type,
					Add:         laboratory1.Laboratory.Information.Address.Add,
				},
				AvgRating: laboratory1.Laboratory.Review.AvgRating,
			})
		}
	}

	if len(laboratoryData) == 0 {
		return c.Status(fiber.StatusOK).JSON(laboratory.GetLaboratoryResponseDto{
			Status:  false,
			Message: "No Laboratory data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(laboratory.GetLaboratoryResponseDto{
		Status:  true,
		Message: "Successfully fetched laboratory data.",
		Data:    laboratoryData,
	})
}
