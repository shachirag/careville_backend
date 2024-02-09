package hospitals

import (
	"careville_backend/database"
	hospitals "careville_backend/dto/customer/hospitals"
	"careville_backend/entity"
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

// @Summary Get Hospitals
// @Tags customer hospitals
// @Description Get Hospitals
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param search query string false "Filter hospitals by search"
// @Param long query float64 false "Longitude for memories sorting (required for distance sorting)"
// @Param lat query float64 false "Latitude for memories sorting (required for distance sorting)"
// @Produce json
// @Success 200 {object} hospitals.GetHospitalResDto
// @Router /customer/healthFacility/get-hospitals [get]
func GetHospitals(c *fiber.Ctx) error {

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
			return c.Status(fiber.StatusBadRequest).JSON(hospitals.GetHospitalResDto{
				Status:  false,
				Message: "Invalid latitude format",
			})
		}

		long, err = strconv.ParseFloat(longParam, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(hospitals.GetHospitalResDto{
				Status:  false,
				Message: "Invalid longitude format",
			})
		}
	}

	filter := bson.M{
		"role":                 "healthFacility",
		"facilityOrProfession": "hospClinic",
	}

	maxDistance := 50000

	if latParam != "" && longParam != "" {
		filter["hospClinic.information.address"] = bson.M{
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
		filter["hospClinic.information.name"] = bson.M{"$regex": searchTitle, "$options": "i"}
	}

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	projection := bson.M{
		"hospClinic.information.name":  1,
		"hospClinic.information.image": 1,
		"_id":                          1,
		"hospClinic.review.avgRating":  1,
		"hospClinic.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
	}

	findOptions := options.Find().SetProjection(projection)

	cursor, err := serviceColl.Find(ctx, filter, sortOptions, findOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(hospitals.GetHospitalResDto{
			Status:  false,
			Message: "Failed to fetch hospitals data: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var hospitalData []hospitals.GetHospitalsRes
	for cursor.Next(ctx) {
		var hospital entity.ServiceEntity
		if err := cursor.Decode(&hospital); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(hospitals.GetHospitalResDto{
				Status:  false,
				Message: "Failed to decode hospitals data: " + err.Error(),
			})
		}
		if hospital.HospClinic != nil {
			hospitalData = append(hospitalData, hospitals.GetHospitalsRes{
				Id:    hospital.Id,
				Image: hospital.HospClinic.Information.Image,
				Name:  hospital.HospClinic.Information.Name,
				Address: hospitals.Address{
					Coordinates: hospital.HospClinic.Information.Address.Coordinates,
					Type:        hospital.HospClinic.Information.Address.Type,
					Add:         hospital.HospClinic.Information.Address.Add,
				},
				AvgRating: hospital.HospClinic.Review.AvgRating,
			})
		}
	}

	if len(hospitalData) == 0 {
		return c.Status(fiber.StatusOK).JSON(hospitals.GetHospitalResDto{
			Status:  false,
			Message: "No Hospital data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(hospitals.GetHospitalResDto{
		Status:  true,
		Message: "Successfully fetched hospitals data.",
		Data:    hospitalData,
	})
}
