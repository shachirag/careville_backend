package services

import (
	"careville_backend/database"
	services "careville_backend/dto/provider/services"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// @Summary Fetch All misc data
// @Description Fetch All misc data
// @Tags services
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} services.MiscResDto
// @Router /provider/services/get-misc-data [get]
func FetchAllMiscData(c *fiber.Ctx) error {
	miscColl := database.GetCollection("misc")

	var healthFacility entity.MiscEntity
	healthFacilityFilter := bson.M{"_id": "healthFacility"}
	err := miscColl.FindOne(ctx, healthFacilityFilter).Decode(&healthFacility)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.MiscResDto{
			Status:  false,
			Message: "Failed to fetch misc data: " + err.Error(),
		})
	}

	response := services.MiscResDto{
		Status:  true,
		Message: "Miscellaneous data fetched successfully",
		Data: services.MiscRes{
			HospClinic: services.HospClinicEntity{
				OtherServices: healthFacility.HospClinic.OtherServices,
				Insurances:    healthFacility.HospClinic.Insurances,
			},
			Laboratory: services.LaboratoryEntity{
				Investigations: healthFacility.Laboratory.Investigations,
			},
			FitnessCenter: services.FitnessCenterEntity{
				Categories: healthFacility.FitnessCenter.Categories,
			},
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
