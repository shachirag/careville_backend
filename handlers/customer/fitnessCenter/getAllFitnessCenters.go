package fitnessCenter

import (
	"careville_backend/database"
	fitnessCenter "careville_backend/dto/customer/fitnessCenter"
	"careville_backend/entity"
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

// @Summary Get fitnessCenter
// @Tags customer fitnessCenter
// @Description Get fitnessCenter
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param search query string false "Filter fitnessCenter by search"
// @Param long query float64 false "Longitude for memories sorting (required for distance sorting)"
// @Param lat query float64 false "Latitude for memories sorting (required for distance sorting)"
// @Produce json
// @Success 200 {object} fitnessCenter.GetFitnessCenterResponseDto
// @Router /customer/healthFacility/get-fitnessCenters [get]
func GetFitnessCenter(c *fiber.Ctx) error {

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
			return c.Status(fiber.StatusBadRequest).JSON(fitnessCenter.GetFitnessCenterResponseDto{
				Status:  false,
				Message: "Invalid latitude format",
			})
		}

		long, err = strconv.ParseFloat(longParam, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fitnessCenter.GetFitnessCenterResponseDto{
				Status:  false,
				Message: "Invalid longitude format",
			})
		}
	}

	filter := bson.M{
		"role":                 "healthFacility",
		"facilityOrProfession": "fitnessCenter",
	}

	maxDistance := 50000

	if latParam != "" && longParam != "" {
		filter["fitnessCenter.information.address"] = bson.M{
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
		filter["fitnessCenter.information.name"] = bson.M{"$regex": searchTitle, "$options": "i"}
	}

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	projection := bson.M{
		"fitnessCenter.information.name":  1,
		"fitnessCenter.information.image": 1,
		"_id":                             1,
		"fitnessCenter.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
		"avgRating": 1,
	}

	findOptions := options.Find().SetProjection(projection)

	cursor, err := serviceColl.Find(ctx, filter, sortOptions, findOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fitnessCenter.GetFitnessCenterResponseDto{
			Status:  false,
			Message: "Failed to fetch fitnessCenter data: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var fitnessCenterData []fitnessCenter.GetFitnessCenterRes
	for cursor.Next(ctx) {
		var fitnessCenter1 entity.ServiceEntity
		if err := cursor.Decode(&fitnessCenter1); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fitnessCenter.GetFitnessCenterResponseDto{
				Status:  false,
				Message: "Failed to decode fitnessCenter data: " + err.Error(),
			})
		}
		if fitnessCenter1.FitnessCenter != nil {
			fitnessCenterData = append(fitnessCenterData, fitnessCenter.GetFitnessCenterRes{
				Id:    fitnessCenter1.Id,
				Image: fitnessCenter1.FitnessCenter.Information.Image,
				Name:  fitnessCenter1.FitnessCenter.Information.Name,
				Address: fitnessCenter.Address{
					Coordinates: fitnessCenter1.FitnessCenter.Information.Address.Coordinates,
					Type:        fitnessCenter1.FitnessCenter.Information.Address.Type,
					Add:         fitnessCenter1.FitnessCenter.Information.Address.Add,
				},
				AvgRating: fitnessCenter1.AvgRating,
			})
		}
	}

	if len(fitnessCenterData) == 0 {
		return c.Status(fiber.StatusOK).JSON(fitnessCenter.GetFitnessCenterResponseDto{
			Status:  false,
			Message: "No FitnessCenter data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fitnessCenter.GetFitnessCenterResponseDto{
		Status:  true,
		Message: "Successfully fetched fitnessCenters data.",
		Data:    fitnessCenterData,
	})
}
