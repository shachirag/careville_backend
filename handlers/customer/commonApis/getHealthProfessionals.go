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
			return c.Status(fiber.StatusBadRequest).JSON(common.GetHealthProfessionalResDto{
				Status:  false,
				Message: "Invalid latitude format",
			})
		}

		long, err = strconv.ParseFloat(longParam, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(common.GetHealthProfessionalResDto{
				Status:  false,
				Message: "Invalid longitude format",
			})
		}
	}

	filter := bson.M{
		"role": "healthProfessional",
	}

	var maxDistance int64
	switch c.Query("facilityOrProfession") {
	case "doctor":
		maxDistance = 50000
		filter["facilityOrProfession"] = "doctor"
	case "physiotherapist":
		maxDistance = 30000
		filter["facilityOrProfession"] = "physiotherapist"
	case "nurse":
		maxDistance = 40000
		filter["facilityOrProfession"] = "nurse"
	case "medicalLabScientist":
		maxDistance = 60000
		filter["facilityOrProfession"] = "medicalLabScientist"
	default:
		maxDistance = 50000
	}

	if latParam != "" && longParam != "" {
		fieldName := "doctor.information.address"
		if filter["facilityOrProfession"] != "doctor" {
			fieldName = "hospClinic.information.address"
		} else if filter["facilityOrProfession"] != "physiotherapist" {
			fieldName = "physiotherapist.information.address"
		} else if filter["facilityOrProfession"] != "nurse" {
			fieldName = "nurse.information.address"
		} else if filter["facilityOrProfession"] != "medicalLabScientist" {
			fieldName = "medicalLabScientist.information.address"
		}
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

	if searchTitle != "" {
		fieldName := "doctor.information.name"
		if filter["facilityOrProfession"] != "doctor" {
			fieldName = "doctor.information.name"
		} else if filter["facilityOrProfession"] != "physiotherapist" {
			fieldName = "physiotherapist.information.name"
		} else if filter["facilityOrProfession"] != "nurse" {
			fieldName = "nurse.information.name"
		} else if filter["facilityOrProfession"] != "medicalLabScientist" {
			fieldName = "medicalLabScientist.information.name"
		}
		filter[fieldName] = bson.M{"$regex": searchTitle, "$options": "i"}
	}

	limit := int64(5)

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1}).SetLimit(limit)

	cursor, err := serviceColl.Find(ctx, filter, sortOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.GetHealthProfessionalResDto{
			Status:  false,
			Message: "Failed to fetch health professionals data: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var healthProfessionalData common.HealthProfessionalResDto
	for cursor.Next(ctx) {
		var healthProfessional entity.ServiceEntity
		if err := cursor.Decode(&healthProfessional); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(common.GetHealthProfessionalResDto{
				Status:  false,
				Message: "Failed to decode health professionals data: " + err.Error(),
			})
		}

		switch healthProfessional.FacilityOrProfession {
		case "doctor":
			if healthProfessional.Doctor != nil {
				healthProfessionalData.Doctors = append(healthProfessionalData.Doctors, common.GetDoctorHealthProfessionalRes{
					Id:         healthProfessional.Id,
					Image:      healthProfessional.Doctor.Information.Image,
					Name:       healthProfessional.Doctor.Information.Name,
					AvgRating:  healthProfessional.Doctor.Review.AvgRating,
					Speciality: healthProfessional.Doctor.AdditionalServices.Speciality,
				})
			}
		case "physiotherapist":
			if healthProfessional.Physiotherapist != nil {
				healthProfessionalData.Physiotherapists = append(healthProfessionalData.Physiotherapists, common.GetHealthProfessionalRes{
					Id:        healthProfessional.Id,
					Image:     healthProfessional.Physiotherapist.Information.Image,
					Name:      healthProfessional.Physiotherapist.Information.Name,
					AvgRating: healthProfessional.Physiotherapist.Review.AvgRating,
				})
			}
		case "nurse":
			if healthProfessional.Nurse != nil {
				healthProfessionalData.Nurse = append(healthProfessionalData.Nurse, common.GetHealthProfessionalRes{
					Id:        healthProfessional.Id,
					Image:     healthProfessional.Nurse.Information.Image,
					Name:      healthProfessional.Nurse.Information.Name,
					AvgRating: healthProfessional.Nurse.Review.AvgRating,
				})
			}
		case "medicalLabScientist":
			if healthProfessional.MedicalLabScientist != nil {
				healthProfessionalData.MedicalLabScientists = append(healthProfessionalData.MedicalLabScientists, common.GetHealthProfessionalRes{
					Id:        healthProfessional.Id,
					Image:     healthProfessional.MedicalLabScientist.Information.Image,
					Name:      healthProfessional.MedicalLabScientist.Information.Name,
					AvgRating: healthProfessional.MedicalLabScientist.Review.AvgRating,
				})
			}
		}
	}

	if len(healthProfessionalData.Physiotherapists) == 0 {
		healthProfessionalData.Physiotherapists = []common.GetHealthProfessionalRes{}
	}
	if len(healthProfessionalData.Doctors) == 0 {
		healthProfessionalData.Doctors = []common.GetDoctorHealthProfessionalRes{}
	}
	if len(healthProfessionalData.Nurse) == 0 {
		healthProfessionalData.Nurse = []common.GetHealthProfessionalRes{}
	}
	if len(healthProfessionalData.MedicalLabScientists) == 0 {
		healthProfessionalData.MedicalLabScientists = []common.GetHealthProfessionalRes{}
	}

	return c.Status(fiber.StatusOK).JSON(common.GetHealthProfessionalResDto{
		Status:  true,
		Message: "Successfully fetched health professionals data.",
		Data: common.HealthProfessionalResDto{
			Doctors:              healthProfessionalData.Doctors,
			Physiotherapists:     healthProfessionalData.Physiotherapists,
			Nurse:                healthProfessionalData.Nurse,
			MedicalLabScientists: healthProfessionalData.MedicalLabScientists,
		},
	})
}
