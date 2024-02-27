package pharmacy

import (
	"careville_backend/database"
	pharmacy "careville_backend/dto/customer/pharmacy"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get pharmacy by ID
// @Tags customer pharmacy
// @Description Get pharmacy by ID
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "pharmacy ID"
// @Produce json
// @Success 200 {object} pharmacy.GetPharmacyResDto
// @Router /customer/healthFacility/get-pharmacy/{id} [get]
func GetPharmacyByID(c *fiber.Ctx) error {
	var (
		serviceColl = database.GetCollection("service")
	)

	idParam := c.Params("id")
	pharmacyID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pharmacy.GetPharmacyResDto{
			Status:  false,
			Message: "Invalid pharmacy ID",
		})
	}

	filter := bson.M{"_id": pharmacyID}

	projection := bson.M{
		"pharmacy.information.name":               1,
		"pharmacy.information.image":              1,
		"_id":                                     1,
		"pharmacy.information.additionalText":     1,
		"pharmacy.review.totalReviews":            1,
		"pharmacy.review.avgRating":               1,
		"pharmacy.additionalServices.id":          1,
		"pharmacy.additionalServices.name":        1,
		"pharmacy.additionalServices.information": 1,
		"pharmacy.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
	}

	findOptions := options.FindOne().SetProjection(projection)

	var pharmacyData entity.ServiceEntity
	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&pharmacyData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(pharmacy.GetPharmacyResDto{
			Status:  false,
			Message: "Failed to fetch pharmacy data: " + err.Error(),
		})
	}

	if pharmacyData.Pharmacy == nil {
		return c.Status(fiber.StatusNotFound).JSON(pharmacy.GetPharmacyResDto{
			Status:  false,
			Message: "pharmacy data not found",
		})
	}

	servicesData := make([]pharmacy.AdditionalServices, 0)
	if pharmacyData.Pharmacy != nil && len(pharmacyData.Pharmacy.AdditionalServices) > 0 {
		for _, service := range pharmacyData.Pharmacy.AdditionalServices {
			servicesData = append(servicesData, pharmacy.AdditionalServices{
				Id:          service.Id,
				Name:        service.Name,
				Information: service.Information,
			})
		}
	}

	laboratoryRes := pharmacy.GetPharmacyResDto{
		Status:  true,
		Message: "Pharmacy data fetched successfully",
		Data: pharmacy.PharmacyResponse{
			Id:                 pharmacyData.Id,
			Image:              pharmacyData.Pharmacy.Information.Image,
			Name:               pharmacyData.Pharmacy.Information.Name,
			AboutUs:            pharmacyData.Pharmacy.Information.AdditionalText,
			Address:            pharmacy.Address(pharmacyData.Pharmacy.Information.Address),
			AdditionalServices: servicesData,
			TotalReviews:       pharmacyData.Pharmacy.Review.TotalReviews,
			AvgRating:          pharmacyData.Pharmacy.Review.AvgRating,
		},
	}

	return c.Status(fiber.StatusOK).JSON(laboratoryRes)
}
