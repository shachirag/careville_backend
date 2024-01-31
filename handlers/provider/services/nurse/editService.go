package nurse

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Update other service
// @Description Update other service
// @Tags nurse
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param serviceId path string true "service ID"
// @Param provider body services.UpdateNurseServiceReqDto true "Update data of service"
// @Produce json
// @Success 200 {object} services.DeleteNurseProfessionalInfoResDto
// @Router /provider/services/update-nurse-service/{serviceId} [put]
func UpdateNurseServiceInfo(c *fiber.Ctx) error {

	var (
		serviceColl = database.GetCollection("service")
		data        services.UpdateNurseServiceReqDto
		provider    entity.ServiceEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(services.DeleteNurseProfessionalInfoResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	serviceId := c.Params("serviceId")
	serviceObjID, err := primitive.ObjectIDFromHex(serviceId)

	if err != nil {
		return c.Status(400).JSON(services.UpdateDoctorImageResDto{
			Status:  false,
			Message: "invalid objectId " + err.Error(),
		})
	}

	filter := bson.M{
		"_id": providerData.ProviderId,
		"nurse.schedule": bson.M{
			"$elemMatch": bson.M{
				"id": serviceObjID,
			},
		},
	}

	err = serviceColl.FindOne(ctx, filter).Decode(&provider)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.DeleteNurseProfessionalInfoResDto{
				Status:  false,
				Message: "service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.DeleteNurseProfessionalInfoResDto{
			Status:  false,
			Message: "Failed to fetch service from MongoDB: " + err.Error(),
		})
	}

	update := bson.M{
		"$set": bson.M{
			"nurse.schedule.$.serviceFees": data.ServiceFees,
			"nurse.schedule.$.name":        data.Name,
			"nurse.schedule.$.slots":       bson.A{},
			"updatedAt":                    time.Now().UTC(),
		},
	}

	// Clearing existing schedule
	// update["$set"].(bson.M)["hospClinic.doctor.$.schedule"] = bson.A{}

	for _, slot := range data.Slots {
		slotUpdate := bson.M{
			"startTime": slot.StartTime,
			"endTime":   slot.EndTime,
			"days":      slot.Days,
		}
		// fmt.Print(schedule)
		update["$set"].(bson.M)["nurse.schedule.$.slots"] = append(update["$set"].(bson.M)["nurse.schedule.$.slots"].(bson.A), slotUpdate)
	}

	// Execute the update operation
	updateRes, err := serviceColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.DeleteNurseProfessionalInfoResDto{
			Status:  false,
			Message: "Failed to update service data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(services.DeleteNurseProfessionalInfoResDto{
			Status:  false,
			Message: "service not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(services.DeleteNurseProfessionalInfoResDto{
		Status:  true,
		Message: "Service data updated successfully",
	})
}
