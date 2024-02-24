package common

import (
	"careville_backend/database"
	common "careville_backend/dto/customer/commonApis"
	"careville_backend/entity"
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

// @Summary Get nearest health facilties
// @Tags customer commonApis
// @Description Get nearest health facilties
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param search query string false "Filter hospitals by search"
// @Param long query float64 false "Longitude for memories sorting (required for distance sorting)"
// @Param lat query float64 false "Latitude for memories sorting (required for distance sorting)"
// @Produce json
// @Success 200 {object} common.GetHealthFacilityResDto
// @Router /customer/healthFacility/get-health-facilities [get]
func GetHealthFacilties(c *fiber.Ctx) error {

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
			return c.Status(fiber.StatusBadRequest).JSON(common.GetHealthFacilityResDto{
				Status:  false,
				Message: "Invalid latitude format",
			})
		}

		long, err = strconv.ParseFloat(longParam, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(common.GetHealthFacilityResDto{
				Status:  false,
				Message: "Invalid longitude format",
			})
		}
	}

	filter := bson.M{
		"role": "healthFacility",
	}

	// if latParam != "" && longParam != "" {
	// 	filter["$or"] = []bson.M{
	// 		{"hospClinic.information.address": bson.M{"$nearSphere": bson.M{"$geometry": bson.M{"type": "Point", "coordinates": []float64{long, lat}}, "$maxDistance": 50000}}},
	// 		{"laboratory.information.address": bson.M{"$nearSphere": bson.M{"$geometry": bson.M{"type": "Point", "coordinates": []float64{long, lat}}, "$maxDistance": 50000}}},
	// 		{"pharmacy.information.address": bson.M{"$nearSphere": bson.M{"$geometry": bson.M{"type": "Point", "coordinates": []float64{long, lat}}, "$maxDistance": 50000}}},
	// 		{"fitnessCenter.information.address": bson.M{"$nearSphere": bson.M{"$geometry": bson.M{"type": "Point", "coordinates": []float64{long, lat}}, "$maxDistance": 50000}}},
	// 	}
	// }

	var maxDistance int64 = 50000

    addGeoQuery := func(fieldName string) {
        filter[fieldName] = bson.M{
            "$nearSphere": bson.M{
                "$geometry": bson.M{
                    "type":        "Point",
                    "coordinates": []float64{long, lat},
                },
                "$maxDistance": maxDistance,
            },
        }
    }

    addGeoQuery("hospClinic.information.address")
    addGeoQuery("laboratory.information.address")
    addGeoQuery("pharmacy.information.address")
    addGeoQuery("fitnessCenter.information.address")

	// var maxDistance int64
	// switch c.Query("facilityOrProfession") {
	// case "hospClinic":
	// 	maxDistance = 50000
	// 	filter["facilityOrProfession"] = "hospClinic"
	// case "laboratory":
	// 	maxDistance = 50000
	// 	filter["facilityOrProfession"] = "laboratory"
	// case "pharmacy":
	// 	maxDistance = 50000
	// 	filter["facilityOrProfession"] = "pharmacy"
	// case "fitnessCenter":
	// 	maxDistance = 50000
	// 	filter["facilityOrProfession"] = "fitnessCenter"
	// default:
	// 	maxDistance = 50000
	// }

	// if latParam != "" && longParam != "" {
	// 	fieldName := ""
	// 	switch c.Query("facilityOrProfession") {
	// 	case "hospClinic":
	// 		fieldName = "hospClinic.information.address"
	// 	case "laboratory":
	// 		fieldName = "laboratory.information.address"
	// 	case "pharmacy":
	// 		fieldName = "pharmacy.information.address"
	// 	case "fitnessCenter":
	// 		fieldName = "fitnessCenter.information.address"
	// 	default:
	// 		fieldName = "hospClinic.information.address" 
	// 	}
	
	// 	filter[fieldName] = bson.M{
	// 		"$nearSphere": bson.M{
	// 			"$geometry": bson.M{
	// 				"type":        "Point",
	// 				"coordinates": []float64{long, lat},
	// 			},
	// 			"$maxDistance": maxDistance,
	// 		},
	// 	}
	// }
	

	if searchTitle != "" {
		filter["$or"] = []bson.M{
			{"hospClinic.information.name": bson.M{"$regex": searchTitle, "$options": "i"}},
			{"laboratory.information.name": bson.M{"$regex": searchTitle, "$options": "i"}},
			{"pharmacy.information.name": bson.M{"$regex": searchTitle, "$options": "i"}},
			{"fitnessCenter.information.name": bson.M{"$regex": searchTitle, "$options": "i"}},
		}
	}

	limit := int64(5)

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1}).SetLimit(limit)

	cursor, err := serviceColl.Find(ctx, filter, sortOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.GetHealthFacilityResDto{
			Status:  false,
			Message: "Failed to fetch health facilities data: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var healthFacilityData common.HealthFacilityResDto
	for cursor.Next(ctx) {
		var healthFacility entity.ServiceEntity
		if err := cursor.Decode(&healthFacility); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(common.GetHealthFacilityResDto{
				Status:  false,
				Message: "Failed to decode health facilities data: " + err.Error(),
			})
		}

		switch healthFacility.FacilityOrProfession {
		case "hospClinic":
			if healthFacility.HospClinic != nil {
				healthFacilityData.HospitalRes = append(healthFacilityData.HospitalRes, common.GetHealthFacilityRes{
					Id:        healthFacility.Id,
					Image:     healthFacility.HospClinic.Information.Image,
					Name:      healthFacility.HospClinic.Information.Name,
					AvgRating: healthFacility.HospClinic.Review.AvgRating,
				})
			}
		case "laboratory":
			if healthFacility.Laboratory != nil {
				healthFacilityData.LaboratoryRes = append(healthFacilityData.LaboratoryRes, common.GetHealthFacilityRes{
					Id:        healthFacility.Id,
					Image:     healthFacility.Laboratory.Information.Image,
					Name:      healthFacility.Laboratory.Information.Name,
					AvgRating: healthFacility.Laboratory.Review.AvgRating,
				})
			}
		case "pharmacy":
			if healthFacility.Pharmacy != nil {
				healthFacilityData.Pharmacy = append(healthFacilityData.Pharmacy, common.GetHealthFacilityRes{
					Id:        healthFacility.Id,
					Image:     healthFacility.Pharmacy.Information.Image,
					Name:      healthFacility.Pharmacy.Information.Name,
					AvgRating: healthFacility.Pharmacy.Review.AvgRating,
				})
			}
		case "fitnessCenter":
			if healthFacility.FitnessCenter != nil {
				healthFacilityData.FitnessCenters = append(healthFacilityData.FitnessCenters, common.GetHealthFacilityRes{
					Id:        healthFacility.Id,
					Image:     healthFacility.FitnessCenter.Information.Image,
					Name:      healthFacility.FitnessCenter.Information.Name,
					AvgRating: healthFacility.FitnessCenter.Review.AvgRating,
				})
			}
		}
	}

	if len(healthFacilityData.HospitalRes) == 0 {
		healthFacilityData.HospitalRes = []common.GetHealthFacilityRes{}
	}
	if len(healthFacilityData.LaboratoryRes) == 0 {
		healthFacilityData.LaboratoryRes = []common.GetHealthFacilityRes{}
	}
	if len(healthFacilityData.Pharmacy) == 0 {
		healthFacilityData.Pharmacy = []common.GetHealthFacilityRes{}
	}
	if len(healthFacilityData.FitnessCenters) == 0 {
		healthFacilityData.FitnessCenters = []common.GetHealthFacilityRes{}
	}

	return c.Status(fiber.StatusOK).JSON(common.GetHealthFacilityResDto{
		Status:  true,
		Message: "Successfully fetched health facilties data.",
		Data: common.HealthFacilityResDto{
			HospitalRes:    healthFacilityData.HospitalRes,
			LaboratoryRes:  healthFacilityData.LaboratoryRes,
			Pharmacy:       healthFacilityData.Pharmacy,
			FitnessCenters: healthFacilityData.FitnessCenters,
		},
	})
}
