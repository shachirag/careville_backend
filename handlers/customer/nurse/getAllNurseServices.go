package nurse

import (
	"careville_backend/database"
	"careville_backend/dto/customer/nurse"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get all services
// @Description Get all services
// @Tags customer nurse
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param nurseId query string true "nurse ID"
// @Produce json
// @Success 200 {object} nurse.NurseServicesResDto
// @Router /customer/healthProfessional/get-nurse-services [get]
func GetNurseServices(c *fiber.Ctx) error {

	var service entity.ServiceEntity

	idParam := c.Query("nurseId")
	nurseID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(nurse.GetNurseResDto{
			Status:  false,
			Message: "Invalid nurse ID",
		})
	}

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"_id": nurseID,
	}

	projection := bson.M{
		"nurse.schedule.id":          1,
		"nurse.schedule.name":        1,
		"nurse.schedule.serviceFees": 1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(nurse.NurseServicesResDto{
				Status:  false,
				Message: "service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(nurse.NurseServicesResDto{
			Status:  false,
			Message: "Failed to fetch services from MongoDB: " + err.Error(),
		})
	}

	serviceData := make([]nurse.NurseServiceRes, 0)
	if service.Nurse != nil && len(service.Nurse.Schedule) > 0 {
		for _, investigation := range service.Nurse.Schedule {
			serviceData = append(serviceData, nurse.NurseServiceRes{
				Id:          investigation.Id,
				Name:        investigation.Name,
				ServiceFees: investigation.ServiceFees,
			})
		}
	}

	if len(serviceData) == 0 {
		return c.Status(fiber.StatusOK).JSON(nurse.NurseServicesResDto{
			Status:  false,
			Message: "No service data found.",
		})
	}

	return c.Status(fiber.StatusOK).JSON(nurse.NurseServicesResDto{
		Status:  true,
		Message: "services retrieved successfully",
		Data:    serviceData,
	})
}
