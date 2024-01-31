package doctorProfession

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

// @Summary Update slot
// @Description Update slot
// @Tags doctorProfession
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param slotId path string true "slot ID"
// @Param provider body services.UpdateDoctorProfessionSlotReqDto true "Update data of service"
// @Produce json
// @Success 200 {object} services.UpdateDoctorProfessionSlotReqDto
// @Router /provider/services/update-doctorProfession-slot/{slotId} [put]
func UpdateDoctorProfessionSlot(c *fiber.Ctx) error {

	var (
		serviceColl = database.GetCollection("service")
		data        services.EditDoctorProfessionSlotReqDto
		provider    entity.ServiceEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(services.UpdateDoctorProfessionSlotResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	slotId := c.Params("slotId")
	slotObjID, err := primitive.ObjectIDFromHex(slotId)

	if err != nil {
		return c.Status(400).JSON(services.UpdateDoctorImageResDto{
			Status:  false,
			Message: "invalid objectId " + err.Error(),
		})
	}

	filter := bson.M{
		"_id": providerData.ProviderId,
		"doctor.schedule.slots": bson.M{
			"$elemMatch": bson.M{
				"id": slotObjID,
			},
		},
	}

	err = serviceColl.FindOne(ctx, filter).Decode(&provider)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.UpdateDoctorProfessionSlotResDto{
				Status:  false,
				Message: "slot not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateDoctorProfessionSlotResDto{
			Status:  false,
			Message: "Failed to fetch slot from MongoDB: " + err.Error(),
		})
	}

	update := bson.M{
		"$set": bson.M{
			"doctor.schedule.slots.$.startTime": data.StartTime,
			"doctor.schedule.slots.$.endTime":   data.EndTime,
			"doctor.schedule.slots.$.days":      data.Days,
			"updatedAt":                         time.Now().UTC(),
		},
	}

	// Execute the update operation
	updateRes, err := serviceColl.UpdateOne(ctx, filter, update)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateDoctorProfessionSlotResDto{
			Status:  false,
			Message: "Failed to update service data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(services.UpdateDoctorProfessionSlotResDto{
			Status:  false,
			Message: "slot not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(services.UpdateDoctorProfessionSlotResDto{
		Status:  true,
		Message: "Slot updated successfully",
	})
}
