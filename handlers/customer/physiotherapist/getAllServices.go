package physiotherapist

import (
	"careville_backend/database"
	"careville_backend/dto/customer/physiotherapist"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get all services
// @Description Get all investigations
// @Tags customer physiotherapist
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param physiotherapistId query string true "physiotherapistId ID"
// @Produce json
// @Success 200 {object} physiotherapist.PhysiotherapistServicesResDto
// @Router /customer/healthProfessional/get-physiotherapist-services [get]
func GetPhysiotherapistServices(c *fiber.Ctx) error {

	var service entity.ServiceEntity

	idParam := c.Query("physiotherapistId")
	physiotherapistID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(physiotherapist.PhysiotherapistServicesResDto{
			Status:  false,
			Message: "Invalid nurse ID",
		})
	}

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"_id": physiotherapistID,
	}

	projection := bson.M{
		"physiotherapist.serviceAndSchedule.id":          1,
		"physiotherapist.serviceAndSchedule.name":        1,
		"physiotherapist.serviceAndSchedule.serviceFees": 1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(physiotherapist.PhysiotherapistServicesResDto{
				Status:  false,
				Message: "service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(physiotherapist.PhysiotherapistServicesResDto{
			Status:  false,
			Message: "Failed to fetch services from MongoDB: " + err.Error(),
		})
	}

	serviceData := make([]physiotherapist.PhysiotherapistServiceRes, 0)
	if service.Physiotherapist != nil && len(service.Physiotherapist.ServiceAndSchedule) > 0 {
		for _, investigation := range service.Physiotherapist.ServiceAndSchedule {
			serviceData = append(serviceData, physiotherapist.PhysiotherapistServiceRes{
				Id:          investigation.Id,
				Name:        investigation.Name,
				ServiceFees: investigation.ServiceFees,
			})
		}
	}

	if len(serviceData) == 0 {
		return c.Status(fiber.StatusOK).JSON(physiotherapist.PhysiotherapistServicesResDto{
			Status:  false,
			Message: "No service data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(physiotherapist.PhysiotherapistServicesResDto{
		Status:  true,
		Message: "services retrieved successfully",
		Data:    serviceData,
	})
}
