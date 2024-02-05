package hospitals

import (
	"careville_backend/database"
	"careville_backend/dto/customer/hospitals"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get hospital by ID
// @Tags customer hospitals
// @Description Get hospital by ID
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "hospital ID"
// @Produce json
// @Success 200 {object} hospitals.GetHospitalsResDto
// @Router /customer/healthFacility/get-hospital/{id} [get]
func GetHospitalByID(c *fiber.Ctx) error {
	var (
		serviceColl = database.GetCollection("service")
	)

	idParam := c.Params("id")
	hospitalID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(hospitals.GetHospitalsResDto{
			Status:  false,
			Message: "Invalid hospital ID",
		})
	}

	filter := bson.M{"_id": hospitalID}

	projection := bson.M{
		"hospClinic.information.name":           1,
		"hospClinic.information.image":          1,
		"hospClinic.information.id":             1,
		"totalReviews":                          1,
		"avgRating":                             1,
		"hospClinic.information.additionalText": 1,
		"hospClinic.otherServices":              1,
		"hospClinic.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
	}

	findOptions := options.FindOne().SetProjection(projection)

	var hospitalData entity.ServiceEntity
	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&hospitalData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(hospitals.GetHospitalsResDto{
			Status:  false,
			Message: "Failed to fetch hospital data: " + err.Error(),
		})
	}

	if hospitalData.HospClinic == nil {
		return c.Status(fiber.StatusNotFound).JSON(hospitals.GetHospitalsResDto{
			Status:  false,
			Message: "Hospital data not found",
		})
	}

	hospitalRes := hospitals.GetHospitalsResDto{
		Status:  true,
		Message: "Hospital data fetched successfully",
		Data: hospitals.HospitalsResponse{
			Id:            hospitalData.Id,
			Image:         hospitalData.HospClinic.Information.Image,
			Name:          hospitalData.HospClinic.Information.Name,
			AboutUs:       hospitalData.HospClinic.Information.AdditionalText,
			Address:       hospitals.Address(hospitalData.HospClinic.Information.Address),
			OtherServices: hospitalData.HospClinic.OtherServices,
			TotalReviews:  hospitalData.TotalReviews,
			AvgRating:     hospitalData.AvgRating,
		},
	}

	return c.Status(fiber.StatusOK).JSON(hospitalRes)
}
