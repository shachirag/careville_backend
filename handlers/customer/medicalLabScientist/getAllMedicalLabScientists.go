package medicalLabScientist

import (
	"careville_backend/database"
	medicalLabScientist "careville_backend/dto/customer/medicalLabScientist"
	"careville_backend/entity"
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

// @Summary Get medicalLabScientist
// @Tags customer medicalLabScientist
// @Description Get medicalLabScientist
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param search query string false "Filter medicalLabScientist by search"
// @Param long query float64 false "Longitude for memories sorting (required for distance sorting)"
// @Param lat query float64 false "Latitude for memories sorting (required for distance sorting)"
// @Produce json
// @Success 200 {object} medicalLabScientist.GetMedicalLabScientistResponseDto
// @Router /customer/healthProfessional/get-medicalLabScientists [get]
func GetMedicalLabScientist1(c *fiber.Ctx) error {

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
			return c.Status(fiber.StatusBadRequest).JSON(medicalLabScientist.GetMedicalLabScientistResponseDto{
				Status:  false,
				Message: "Invalid latitude format",
			})
		}

		long, err = strconv.ParseFloat(longParam, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(medicalLabScientist.GetMedicalLabScientistResponseDto{
				Status:  false,
				Message: "Invalid longitude format",
			})
		}
	}

	filter := bson.M{
		"role":                 "healthProfessional",
		"facilityOrProfession": "medicalLabScientist",
	}

	maxDistance := 50000

	if latParam != "" && longParam != "" {
		filter["medicalLabScientist.information.address"] = bson.M{
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
		filter["medicalLabScientist.information.name"] = bson.M{"$regex": searchTitle, "$options": "i"}
	}

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	projection := bson.M{
		"medicalLabScientist.information.name":  1,
		"medicalLabScientist.information.image": 1,
		"_id":                                   1,
	}

	findOptions := options.Find().SetProjection(projection)

	cursor, err := serviceColl.Find(ctx, filter, sortOptions, findOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.GetMedicalLabScientistResponseDto{
			Status:  false,
			Message: "Failed to fetch medicalLabScientist data: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var medicalLabScientistData []medicalLabScientist.GetMedicalLabScientistRes
	for cursor.Next(ctx) {
		var medicalLabScientist1 entity.ServiceEntity
		if err := cursor.Decode(&medicalLabScientist1); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.GetMedicalLabScientistResponseDto{
				Status:  false,
				Message: "Failed to decode medicalLabScientist data: " + err.Error(),
			})
		}
		if medicalLabScientist1.MedicalLabScientist != nil {
			medicalLabScientistData = append(medicalLabScientistData, medicalLabScientist.GetMedicalLabScientistRes{
				Id:          medicalLabScientist1.Id,
				Image:       medicalLabScientist1.MedicalLabScientist.Information.Image,
				Name:        medicalLabScientist1.MedicalLabScientist.Information.Name,
				ServiceType: "MedicalLabScientist",
			})
		}
	}

	if len(medicalLabScientistData) == 0 {
		return c.Status(fiber.StatusOK).JSON(medicalLabScientist.GetMedicalLabScientistResponseDto{
			Status:  false,
			Message: "No MedicalLabScientist data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(medicalLabScientist.GetMedicalLabScientistResponseDto{
		Status:  true,
		Message: "Successfully fetched medicalLabScientists data.",
		Data:    medicalLabScientistData,
	})
}
