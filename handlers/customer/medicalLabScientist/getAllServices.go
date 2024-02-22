package medicalLabScientist

import (
	"careville_backend/database"
	"careville_backend/dto/customer/medicalLabScientist"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get all services
// @Description Get all services
// @Tags customer medicalLabScientist
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param medicalLabScientistId query string true "medicalLabScientist ID"
// @Produce json
// @Success 200 {object} medicalLabScientist.MedicalLabScientistServicesResDto
// @Router /customer/healthProfessional/get-medicalLabScientist-services [get]
func GetMedicalLabScientistServices(c *fiber.Ctx) error {

	var service entity.ServiceEntity

	idParam := c.Query("medicalLabScientistId")
	medicalLabScientistID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(medicalLabScientist.MedicalLabScientistServicesResDto{
			Status:  false,
			Message: "Invalid nurse ID",
		})
	}

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"_id": medicalLabScientistID,
	}

	projection := bson.M{
		"medicalLabScientist.serviceAndSchedule.id":          1,
		"medicalLabScientist.serviceAndSchedule.name":        1,
		"medicalLabScientist.serviceAndSchedule.serviceFees": 1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(medicalLabScientist.MedicalLabScientistServicesResDto{
				Status:  false,
				Message: "service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(medicalLabScientist.MedicalLabScientistServicesResDto{
			Status:  false,
			Message: "Failed to fetch services from MongoDB: " + err.Error(),
		})
	}

	serviceData := make([]medicalLabScientist.MedicalLabScientistServiceRes, 0)
	if service.MedicalLabScientist != nil && len(service.MedicalLabScientist.ServiceAndSchedule) > 0 {
		for _, investigation := range service.MedicalLabScientist.ServiceAndSchedule {
			serviceData = append(serviceData, medicalLabScientist.MedicalLabScientistServiceRes{
				Id:          investigation.Id,
				Name:        investigation.Name,
				ServiceFees: investigation.ServiceFees,
			})
		}
	}

	if len(serviceData) == 0 {
		return c.Status(fiber.StatusOK).JSON(medicalLabScientist.MedicalLabScientistServicesResDto{
			Status:  false,
			Message: "No service data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(medicalLabScientist.MedicalLabScientistServicesResDto{
		Status:  true,
		Message: "services retrieved successfully",
		Data:    serviceData,
	})
}
