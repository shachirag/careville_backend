package doctorProfession

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Delete slot
// @Description Delete slot
// @Tags doctorProfession
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param serviceId path string true "serviceId"
// @Produce json
// @Success 200 {object} services.GetPhysiotherapistServicesResDto
// @Router /provider/services/delete-doctorProfession-slot/{slotId} [delete]
func DeleteDoctorsSlots(c *fiber.Ctx) error {
	ctx := context.TODO()

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	slotId := c.Params("slotId")
	slotObjID, err := primitive.ObjectIDFromHex(slotId)
	if err != nil {
		return c.Status(400).JSON(services.GetPhysiotherapistServicesResDto{
			Status:  false,
			Message: "invalid objectId " + err.Error(),
		})
	}

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"_id": providerData.ProviderId,
		"doctor.schedule.slots": bson.M{
			"$elemMatch": bson.M{
				"id": slotObjID,
			},
		},
	}

	// Update to pull the matching service from the array
	update := bson.M{
		"$pull": bson.M{
			"doctor.schedule.slots": bson.M{"id": slotObjID},
		},
	}

	// Perform the update operation
	updateResult, err := serviceColl.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.GetPhysiotherapistServicesResDto{
				Status:  false,
				Message: "slot not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.GetPhysiotherapistServicesResDto{
			Status:  false,
			Message: "Failed to update slot in MongoDB: " + err.Error(),
		})
	}

	if updateResult.ModifiedCount == 0 {
		// If no documents were modified, the service with the given ID was not found
		return c.Status(fiber.StatusNotFound).JSON(services.GetPhysiotherapistServicesResDto{
			Status:  false,
			Message: "No slot information found for the provided service ID",
		})
	}

	return c.Status(fiber.StatusOK).JSON(services.UpdatePhysiotherapistProfessionalInfoResDto{
		Status:  true,
		Message: "Deleted successfully",
	})
}
