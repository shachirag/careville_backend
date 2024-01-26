package hospClinic

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Add other services
// @Tags hospClinic
// @Description Add HospitalClinic
// @Accept multipart/form-data
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param  provider formData services.OtherServiceReqDto true "add HospitalClinic"
// @Produce json
// @Success 200 {object} services.UpdateDoctorResDto
// @Router /provider/services/add-hospClinic-other-services [post]
func AddServices(c *fiber.Ctx) error {
	var (
		servicesColl = database.GetCollection("service")
		data         services.OtherServiceReqDto
		provider     entity.ServiceEntity
	)

	// Parsing the request body
	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(500).JSON(services.UpdateDoctorResDto{
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
			return c.Status(fiber.StatusNotFound).JSON(services.UpdateDoctorImageResDto{
				Status:  false,
				Message: "provider not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.UpdateDoctorImageResDto{
			Status:  false,
			Message: "Failed to fetch provider from MongoDB: " + err.Error(),
		})
	}

	update := bson.M{
		"$pull": bson.M{
			"hospClinic.otherServices": bson.M{
				"$nin": data.OtherServices,
			},
			"hospClinic.insurances": bson.M{
				"$nin": data.Insurances,
			},
		},
	}

	_, err = servicesColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.HospitalClinicResDto{
			Status:  false,
			Message: "Failed to update provider data in MongoDB: " + err.Error(),
		})
	}

	update = bson.M{
		"$addToSet": bson.M{
			"hospClinic.otherServices": bson.M{
				"$each": data.OtherServices,
			}, "hospClinic.insurances": bson.M{
				"$each": data.Insurances,
			},
		},
	}

	updateRes, err := servicesColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(services.HospitalClinicResDto{
			Status:  false,
			Message: "Failed to update provider data in MongoDB: " + err.Error(),
		})
	}

	if updateRes.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(services.HospitalClinicResDto{
			Status:  false,
			Message: "provider not found",
		})
	}

	hospClinicRes := services.HospitalClinicResDto{
		Status:  true,
		Message: "other services and insurances added successfully",
	}
	return c.Status(fiber.StatusOK).JSON(hospClinicRes)
}
