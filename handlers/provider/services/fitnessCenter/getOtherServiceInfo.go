package fitnessCenter

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get other service info
// @Description Get other service info
// @Tags fitnessCenter
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param otherServiceId path string true "other service ID"
// @Produce json
// @Success 200 {object} services.FitnessCenterOtherServiceResDto
// @Router /provider/services/get-fitness-other-service-info/{otherServiceId} [get]
func GetOtherServiceInfo(c *fiber.Ctx) error {

	var service entity.ServiceEntity

	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	otherServiceId := c.Params("otherServiceId")
	otherServiceObjID, err := primitive.ObjectIDFromHex(otherServiceId)

	if err != nil {
		return c.Status(400).JSON(services.FitnessCenterOtherServiceResDto{
			Status:  false,
			Message: "invalid objectId " + err.Error(),
		})
	}

	serviceColl := database.GetCollection("service")

	filter := bson.M{
		"_id": providerData.ProviderId,
		"fitnessCenter.additionalServices": bson.M{
			"$elemMatch": bson.M{
				"id": otherServiceObjID,
			},
		},
	}

	projection := bson.M{
		"fitnessCenter.additionalServices.id":          1,
		"fitnessCenter.additionalServices.name":        1,
		"fitnessCenter.additionalServices.information": 1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.FitnessCenterOtherServiceResDto{
				Status:  false,
				Message: "service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.FitnessCenterOtherServiceResDto{
			Status:  false,
			Message: "Failed to fetch service from MongoDB: " + err.Error(),
		})
	}

	if service.FitnessCenter == nil {
		return c.Status(fiber.StatusNotFound).JSON(services.FitnessCenterOtherServiceResDto{
			Status:  false,
			Message: "No other service information found for the service",
		})
	}

	var servicesRes services.FitnessOtherServiceRes

	for _, services1 := range service.FitnessCenter.AdditionalServices {
		if services1.Id == otherServiceObjID {
			trainerRes := services.FitnessOtherServiceRes{
				Id:          services1.Id,
				Name:        services1.Name,
				Information: services1.Information,
			}

			servicesRes = trainerRes
			break
		}
	}

	return c.Status(fiber.StatusOK).JSON(services.FitnessCenterOtherServiceResDto{
		Status:  true,
		Message: "other service retrieved successfully",
		Data:    servicesRes,
	})
}
