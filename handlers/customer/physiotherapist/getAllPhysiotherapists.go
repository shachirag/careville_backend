package physiotherapist

import (
	"careville_backend/database"
	physiotherapist "careville_backend/dto/customer/physiotherapist"
	"careville_backend/entity"
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

// @Summary Get physiotherapist
// @Tags customer physiotherapist
// @Description Get physiotherapist
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param search query string false "Filter physiotherapist by search"
// @Param long query float64 false "Longitude for memories sorting (required for distance sorting)"
// @Param lat query float64 false "Latitude for memories sorting (required for distance sorting)"
// @Produce json
// @Success 200 {object} physiotherapist.GetPhysiotherapistResponseDto
// @Router /customer/healthProfessional/get-physiotherapists [get]
func GetPhysiotherapists(c *fiber.Ctx) error {

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
			return c.Status(fiber.StatusBadRequest).JSON(physiotherapist.GetPhysiotherapistResponseDto{
				Status:  false,
				Message: "Invalid latitude format",
			})
		}

		long, err = strconv.ParseFloat(longParam, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(physiotherapist.GetPhysiotherapistResponseDto{
				Status:  false,
				Message: "Invalid longitude format",
			})
		}
	}

	filter := bson.M{
		"role":                 "healthProfessional",
		"facilityOrProfession": "physiotherapist",
	}

	maxDistance := 50000

	if latParam != "" && longParam != "" {
		filter["physiotherapist.information.address"] = bson.M{
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
		filter["physiotherapist.information.name"] = bson.M{"$regex": searchTitle, "$options": "i"}
	}

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	projection := bson.M{
		"physiotherapist.information.name":  1,
		"physiotherapist.information.image": 1,
		"_id":                               1,
	}

	findOptions := options.Find().SetProjection(projection)

	cursor, err := serviceColl.Find(ctx, filter, sortOptions, findOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(physiotherapist.GetPhysiotherapistResponseDto{
			Status:  false,
			Message: "Failed to fetch physiotherapist data: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var physiotherapistData []physiotherapist.GetPhysiotherapistRes
	for cursor.Next(ctx) {
		var physiotherapist1 entity.ServiceEntity
		if err := cursor.Decode(&physiotherapist1); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(physiotherapist.GetPhysiotherapistResponseDto{
				Status:  false,
				Message: "Failed to decode physiotherapist data: " + err.Error(),
			})
		}
		if physiotherapist1.Physiotherapist != nil {
			physiotherapistData = append(physiotherapistData, physiotherapist.GetPhysiotherapistRes{
				Id:    physiotherapist1.Id,
				Image: physiotherapist1.Physiotherapist.Information.Image,
				Name:  physiotherapist1.Physiotherapist.Information.Name,
				NextAvailable: physiotherapist.NextAvailable{
					StartTime: "",
					LastTime:  "",
				},
			})
		}
	}

	if len(physiotherapistData) == 0 {
		return c.Status(fiber.StatusOK).JSON(physiotherapist.GetPhysiotherapistResponseDto{
			Status:  false,
			Message: "No Physiotherapist data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(physiotherapist.GetPhysiotherapistResponseDto{
		Status:  true,
		Message: "Successfully fetched physiotherapists data.",
		Data:    physiotherapistData,
	})
}
