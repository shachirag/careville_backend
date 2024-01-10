package services

import (
	"careville_backend/database"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get all doctors
// @Tags services
// @Description Get all doctors
// @Produce json
// @Success 200 {object} services.DoctorResDto
// @Router /provider/get-all-doctors [get]
func GetAllDoctors(c *fiber.Ctx) error {
	var (
		serviceColl = database.GetCollection("service")
	)

	filter := bson.M{"status": "approved"}

	sortOptions := options.Find().SetSort(bson.M{"updatedAt": -1})

	cursor, err := serviceColl.Find(ctx, filter, sortOptions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.DoctorResDto{
			Status:  false,
			Message: "Failed to fetch doctors data: " + err.Error(),
		})
	}
	defer cursor.Close(ctx)

	var doctorData []services.DoctorRes
	for cursor.Next(ctx) {
		var doctor entity.ServiceEntity
		if err := cursor.Decode(&doctor); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(services.DoctorResDto{
				Status:  false,
				Message: "Failed to decode doctor data: " + err.Error(),
			})
		}
		doctorData = append(doctorData, services.DoctorRes{
			// Id: doctor.Id,
		})
	}

	if len(doctorData) == 0 {
		return c.Status(fiber.StatusOK).JSON(services.DoctorResDto{
			Status:  false,
			Message: "No doctor data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(services.DoctorResDto{
		Status:  true,
		Message: "Successfully fetched all doctors",
		Data:    doctorData,
	})
}
