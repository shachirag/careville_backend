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

// @Summary Get nearest health professionals
// @Tags customer commonApis
// @Description Get nearest health professionals
//
// @Param Authorization header string true "Authentication header"
//
// @Param search query string false "Filter health professionals by search"
// @Param long query float64 false "Longitude for memories sorting (required for distance sorting)"
// @Param lat query float64 false "Latitude for memories sorting (required for distance sorting)"
// @Produce json
// @Success 200 {object} common.GetHealthProfessionalResDto
// @Router /customer/healthProfessional/get-health-professionals [get]
func GetHealthProfessionals(c *fiber.Ctx) error {

	searchTitle := c.Query("search", "")

	latParam := c.Query("lat")
	longParam := c.Query("long")

	doctorData, err := getProfessionalsByLocation("doctor", "doctor.information.address", latParam, longParam, "doctor.information.name", searchTitle)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.GetHealthFacilityResDto{
			Status:  false,
			Message: "Failed to get doctor data: " + err.Error(),
		})
	}

	physiotherapistData, err := getProfessionalsByLocation("physiotherapist", "physiotherapist.information.address", latParam, longParam, "physiotherapist.information.name", searchTitle)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.GetHealthFacilityResDto{
			Status:  false,
			Message: "Failed to get physiotherapist data: " + err.Error(),
		})
	}

	nurseData, err := getProfessionalsByLocation("nurse", "nurse.information.address", latParam, longParam, "nurse.information.name", searchTitle)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.GetHealthFacilityResDto{
			Status:  false,
			Message: "Failed to get nurse data: " + err.Error(),
		})
	}

	medicalLabScientistData, err := getProfessionalsByLocation("medicalLabScientist", "medicalLabScientist.information.address", latParam, longParam, "medicalLabScientist.information.name", searchTitle)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.GetHealthFacilityResDto{
			Status:  false,
			Message: "Failed to get medicalLabScientist data: " + err.Error(),
		})
	}

	response := common.GetHealthProfessionalResDto{
		Status:  true,
		Message: "Successfully fetched health professionals data.",
		Data: common.HealthProfessionalResDto{
			Doctors:              []common.GetDoctorHealthProfessionalRes{},
			Physiotherapists:     []common.GetHealthProfessionalRes{},
			Nurse:                []common.GetHealthProfessionalRes{},
			MedicalLabScientists: []common.GetHealthProfessionalRes{},
		},
	}
	for _, entity := range *doctorData {
		switch entity.FacilityOrProfession {
		case "doctor":
			if entity.Doctor != nil {
				response.Data.Doctors = append(response.Data.Doctors, common.GetDoctorHealthProfessionalRes{
					Id:         entity.Id,
					Image:      entity.Doctor.Information.Image,
					Name:       entity.Doctor.Information.Name,
					AvgRating:  entity.Doctor.Review.AvgRating,
					Speciality: entity.Doctor.AdditionalServices.Speciality,
				})
			}
		}
	}

	for _, entity := range *medicalLabScientistData {
		switch entity.FacilityOrProfession {
		case "medicalLabScientist":
			if entity.MedicalLabScientist != nil {
				response.Data.MedicalLabScientists = append(response.Data.MedicalLabScientists, common.GetHealthProfessionalRes{
					Id:        entity.Id,
					Image:     entity.MedicalLabScientist.Information.Image,
					Name:      entity.MedicalLabScientist.Information.Name,
					AvgRating: entity.MedicalLabScientist.Review.AvgRating,
				})
			}
		}
	}

	for _, entity := range *nurseData {
		switch entity.FacilityOrProfession {
		case "nurse":
			if entity.Nurse != nil {
				response.Data.Nurse = append(response.Data.Nurse, common.GetHealthProfessionalRes{
					Id:        entity.Id,
					Image:     entity.Nurse.Information.Image,
					Name:      entity.Nurse.Information.Name,
					AvgRating: entity.Nurse.Review.AvgRating,
				})
			}
		}
	}

	for _, entity := range *physiotherapistData {
		switch entity.FacilityOrProfession {
		case "physiotherapist":
			if entity.Physiotherapist != nil {
				response.Data.Physiotherapists = append(response.Data.Physiotherapists, common.GetHealthProfessionalRes{
					Id:        entity.Id,
					Image:     entity.Physiotherapist.Information.Image,
					Name:      entity.Physiotherapist.Information.Name,
					AvgRating: entity.Physiotherapist.Review.AvgRating,
				})
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func getProfessionalsByLocation(facilityOrProfession string, addressFieldKey string, lat string, lng string, searchFieldKey string, searchQuery string) (*[]entity.ServiceEntity, error) {
	filter := bson.M{
		"role": "healthProfessional",
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
		filter[searchFieldKey] = searchQuery
	}
	limit := int64(5)

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1}).SetLimit(limit)

	cursor, err := database.GetCollection("service").Find(ctx, filter, sortOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var healthProfessionals []entity.ServiceEntity
	err = cursor.All(ctx, &healthProfessionals)
	if err != nil {
		return nil, err
	}
	return &healthProfessionals, nil
}
