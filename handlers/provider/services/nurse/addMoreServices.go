package nurse

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

// @Summary Add more service
// @Tags nurse
// @Description Add more service
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param  provider body services.MoreNurseServiceReqDto true "add AddMoreDoctors"
// @Produce json
// @Success 200 {object} services.MoreNurseServiceResDto
// @Router /provider/services/add-more-nurse-service [post]
func AddMoreNurseServices(c *fiber.Ctx) error {
	var (
		servicesColl = database.GetCollection("service")
		data         services.MoreNurseServiceReqDto
		provider     entity.ServiceEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(services.MoreNurseServiceResDto{
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
			return c.Status(fiber.StatusNotFound).JSON(services.MoreNurseServiceResDto{
				Status:  false,
				Message: "service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.MoreNurseServiceResDto{
			Status:  false,
			Message: "Failed to fetch service from MongoDB: " + err.Error(),
		})
	}

	var slots []entity.Slots
	for _, inv := range data.Slots {
		convertedInv := entity.Slots{
			StartTime: inv.StartTime,
			EndTime:   inv.EndTime,
			Days:      inv.Days,
		}
		slots = append(slots, convertedInv)
	}

	moreService := []entity.ServiceAndSchedule{
		{
			Id:          primitive.NewObjectID(),
			Name:        data.Name,
			ServiceFees: data.ServiceFees,
			Slots:       slots,
		},
	}

	update := bson.M{
		"$push": bson.M{"nurse.schedule": bson.M{"$each": moreService}},
	}

	updateRes, err := servicesColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.HospitalClinicServiceResDto{
			Status:  false,
			Message: "Failed to update service data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(services.HospitalClinicServiceResDto{
			Status:  false,
			Message: "service not found",
		})
	}

	hospClinicRes := services.HospitalClinicServiceResDto{
		Status:  true,
		Message: "Service added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(hospClinicRes)
}
