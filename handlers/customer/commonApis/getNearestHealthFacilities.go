package common

import (
	"careville_backend/database"
	common "careville_backend/dto/customer/commonApis"
	"careville_backend/entity"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
	searchTitle := c.Query("search", "")
	latParam := c.Query("lat")
	longParam := c.Query("long")

	hospitalData, err := getFacilitiesByLocation("hospClinic", "hospClinic.information.address", latParam, longParam, "hospClinic.information.name", searchTitle)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.GetHealthFacilityResDto{
			Status:  false,
			Message: "Failed to get hospital data: " + err.Error(),
		})
	}

	laboratoryData, err := getFacilitiesByLocation("laboratory", "laboratory.information.address", latParam, longParam, "laboratory.information.name", searchTitle)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.GetHealthFacilityResDto{
			Status:  false,
			Message: "Failed to get laboratory data: " + err.Error(),
		})
	}

	pharmacyData, err := getFacilitiesByLocation("pharmacy", "pharmacy.information.address", latParam, longParam, "pharmacy.information.name", searchTitle)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.GetHealthFacilityResDto{
			Status:  false,
			Message: "Failed to get pharmacy data: " + err.Error(),
		})
	}

	fitnessCenterData, err := getFacilitiesByLocation("fitnessCenter", "fitnessCenter.information.address", latParam, longParam, "fitnessCenter.information.name", searchTitle)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.GetHealthFacilityResDto{
			Status:  false,
			Message: "Failed to get fitness center data: " + err.Error(),
		})
	}

	response := common.GetHealthFacilityResDto{
		Status:  true,
		Message: "Successfully fetched health facilities data.",
		Data: common.HealthFacilityResDto{
			HospitalRes:    []common.GetHealthFacilityRes{},
			LaboratoryRes:  []common.GetHealthFacilityRes{},
			Pharmacy:       []common.GetHealthFacilityRes{},
			FitnessCenters: []common.GetHealthFacilityRes{},
		},
	}

	for _, entity := range *hospitalData {
		switch entity.FacilityOrProfession {
		case "hospClinic":
			if entity.HospClinic != nil {
				response.Data.HospitalRes = append(response.Data.HospitalRes, common.GetHealthFacilityRes{
					Id:        entity.Id,
					Image:     entity.HospClinic.Information.Image,
					Name:      entity.HospClinic.Information.Name,
					AvgRating: entity.HospClinic.Review.AvgRating,
				})
			}
		}
	}

	for _, entity := range *laboratoryData {
		switch entity.FacilityOrProfession {
		case "laboratory":
			if entity.Laboratory != nil {
				response.Data.LaboratoryRes = append(response.Data.LaboratoryRes, common.GetHealthFacilityRes{
					Id:        entity.Id,
					Image:     entity.Laboratory.Information.Image,
					Name:      entity.Laboratory.Information.Name,
					AvgRating: entity.Laboratory.Review.AvgRating,
				})
			}
		}
	}

	for _, entity := range *pharmacyData {
		switch entity.FacilityOrProfession {
		case "pharmacy":
			if entity.Pharmacy != nil {
				response.Data.Pharmacy = append(response.Data.Pharmacy, common.GetHealthFacilityRes{
					Id:        entity.Id,
					Image:     entity.Pharmacy.Information.Image,
					Name:      entity.Pharmacy.Information.Name,
					AvgRating: entity.Pharmacy.Review.AvgRating,
				})
			}
		}
	}

	for _, entity := range *fitnessCenterData {
		switch entity.FacilityOrProfession {
		case "fitnessCenter":
			if entity.FitnessCenter != nil {
				response.Data.FitnessCenters = append(response.Data.FitnessCenters, common.GetHealthFacilityRes{
					Id:        entity.Id,
					Image:     entity.FitnessCenter.Information.Image,
					Name:      entity.FitnessCenter.Information.Name,
					AvgRating: entity.FitnessCenter.Review.AvgRating,
				})
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func getFacilitiesByLocation(facilityOrProfession string, addressFieldKey string, lat string, lng string, searchFieldKey string, searchQuery string) (*[]entity.ServiceEntity, error) {
	filter := bson.M{
		"role": "healthFacility",
	}
	filter["facilityOrProfession"] = facilityOrProfession
	if lat != "" && lng != "" {
		lat1, err := strconv.ParseFloat(lat, 64)
		if err != nil {
			return nil, err
		}

		long1, err := strconv.ParseFloat(lng, 64)
		if err != nil {
			return nil, err
		}
		filter[addressFieldKey] = bson.M{
			"$nearSphere": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{long1, lat1},
				},
				"$maxDistance": 50000,
			},
		}
	}
	if searchQuery != "" {
		filter[searchFieldKey] = bson.M{"$regex": searchQuery, "$options": "i"}
	}
	limit := int64(5)

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1}).SetLimit(limit)

	cursor, err := database.GetCollection("service").Find(ctx, filter, sortOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var healthFacilities []entity.ServiceEntity
	err = cursor.All(ctx, &healthFacilities)
	if err != nil {
		return nil, err
	}
	return &healthFacilities, nil
}
