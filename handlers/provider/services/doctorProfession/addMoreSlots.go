package doctorProfession

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Add more slot
// @Tags doctorProfession
// @Description Add more slot
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param  provider body services.UpdateDoctorProfessionSlotReqDto true "add AddMoreDoctors"
// @Produce json
// @Success 200 {object} services.UpdateDoctorProfessionSlotResDto
// @Router /provider/services/add-more-doctorProfession-slot [post]
func AddMoreDoctorProfessionSlots(c *fiber.Ctx) error {
	var (
		servicesColl = database.GetCollection("service")
		data         services.UpdateDoctorProfessionSlotReqDto
		provider     entity.ServiceEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(services.UpdateDoctorProfessionSlotResDto{
			Status:  false,
			Message: err.Error(),
		})
	}

	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	filter := bson.M{
		"_id": providerData.ProviderId,
	}

	err = servicesColl.FindOne(ctx, filter).Decode(&provider)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.UpdateDoctorProfessionSlotResDto{
				Status:  false,
				Message: "service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateDoctorProfessionSlotResDto{
			Status:  false,
			Message: "Failed to fetch service from MongoDB: " + err.Error(),
		})
	}

	// Convert request slots to entity format
	var slots []entity.DoctorSlots
	for _, slot := range data.Slots {
		newSlot := entity.DoctorSlots{
			Id:        primitive.NewObjectID(),
			StartTime: slot.StartTime,
			EndTime:   slot.EndTime,
			Days:      slot.Days,
		}
		slots = append(slots, newSlot)
	}

	// Update the provider's service document to add the new slots
	update := bson.M{
		"$push": bson.M{"doctor.schedule.slots": bson.M{"$each": slots}},
	}

	updateRes, err := servicesColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateDoctorProfessionSlotResDto{
			Status:  false,
			Message: "Failed to update service data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(services.UpdateDoctorProfessionSlotResDto{
			Status:  false,
			Message: "service not found",
		})
	}

	slotRes := services.UpdateDoctorProfessionSlotResDto{
		Status:  true,
		Message: "Service added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(slotRes)
}
