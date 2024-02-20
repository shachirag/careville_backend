package emergency

import (
	"careville_backend/database"
	doctorProfession "careville_backend/dto/customer/doctorProfession"
	"careville_backend/entity"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get emergency doctors
// @Tags emergency
// @Description Get emergency doctors
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param search query string false "Filter doctorProfession by search"
// @Param long query float64 false "Longitude for memories sorting (required for distance sorting)"
// @Param lat query float64 false "Latitude for memories sorting (required for distance sorting)"
// @Produce json
// @Success 200 {object} doctorProfession.EmergencyDoctorResDto
// @Router /customer/emergency/get-emergency-doctors [get]
func GetEmergencyDoctors(c *fiber.Ctx) error {

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
			return c.Status(fiber.StatusBadRequest).JSON(doctorProfession.EmergencyDoctorResDto{
				Status:  false,
				Message: "Invalid latitude format",
			})
		}

		long, err = strconv.ParseFloat(longParam, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(doctorProfession.EmergencyDoctorResDto{
				Status:  false,
				Message: "Invalid longitude format",
			})
		}
	}

	filter := bson.M{
		"role":                 "healthProfessional",
		"facilityOrProfession": "doctor",
		// "serviceStatus":        "approved",
		"doctor.information.isEmergencyAvailable": true,
	}

	maxDistance := 50000

	if latParam != "" && longParam != "" {
		filter["doctor.information.address"] = bson.M{
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
		filter["doctor.information.name"] = bson.M{"$regex": searchTitle, "$options": "i"}
	}

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	projection := bson.M{
		"doctor.information.name":              1,
		"doctor.information.image":             1,
		"_id":                                  1,
		"doctor.additionalServices.speciality": 1,
		"doctor.schedule.consultationFees":     1,
	}

	findOptions := options.Find().SetProjection(projection)

	cursor, err := serviceColl.Find(ctx, filter, sortOptions, findOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(doctorProfession.EmergencyDoctorResDto{
			Status:  false,
			Message: "Failed to fetch doctors data: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var doctorData []doctorProfession.GetEmergencyDoctorProfessionRes
	for cursor.Next(ctx) {
		var doctor entity.ServiceEntity
		if err := cursor.Decode(&doctor); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(doctorProfession.EmergencyDoctorResDto{
				Status:  false,
				Message: "Failed to decode doctors data: " + err.Error(),
			})
		}
		if doctor.Doctor != nil {
			doctorData = append(doctorData, doctorProfession.GetEmergencyDoctorProfessionRes{
				Id:                   doctor.Id,
				Image:                doctor.Doctor.Information.Image,
				Name:                 doctor.Doctor.Information.Name,
				Speciality:           doctor.Doctor.AdditionalServices.Speciality,
				IsEmergencyAvailable: doctor.Doctor.Information.IsEmergencyAvailable,
				ConsultationFees:     doctor.Doctor.Schedule.ConsultationFees,
			})
		}
	}

	if len(doctorData) == 0 {
		return c.Status(fiber.StatusOK).JSON(doctorProfession.EmergencyDoctorResDto{
			Status:  false,
			Message: "No Doctor data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(doctorProfession.EmergencyDoctorResDto{
		Status:  true,
		Message: "Successfully fetched doctors data.",
		Data:    doctorData,
	})
}
