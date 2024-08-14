package fitnessCenter

import (
	"careville_backend/database"
	"careville_backend/dto/customer/fitnessCenter"
	customerMiddleware "careville_backend/dto/customer/middleware"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Get fitnessCenter by ID
// @Tags customer fitnessCenter
// @Description Get fitnessCenter by ID
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Param id path string true "fitnessCenter ID"
// @Produce json
// @Success 200 {object} fitnessCenter.GetFitnessCenterResDto
// @Router /customer/healthFacility/get-fitnessCenter/{id} [get]
func GetFitnessCenterByID(c *fiber.Ctx) error {
	var (
		serviceColl = database.GetCollection("service")
	)

	idParam := c.Params("id")
	fitnessCenterID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fitnessCenter.GetFitnessCenterResDto{
			Status:  false,
			Message: "Invalid fitnessCenter ID",
		})
	}

	customerData := customerMiddleware.GetCustomerMiddlewareData(c)
	customerFilter := bson.M{
		"_id": customerData.CustomerId,
	}

	var customer entity.CustomerEntity
	err = database.GetCollection("customer").FindOne(ctx, customerFilter).Decode(&customer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fitnessCenter.GetFitnessCenterResDto{
			Status:  false,
			Message: "Failed to fetch customer data: " + err.Error(),
		})
	}

	isCustomerFamilyMember := false
	if len(customer.FamilyMembers) > 0 {
		isCustomerFamilyMember = true
	}

	filter := bson.M{"_id": fitnessCenterID}

	projection := bson.M{
		"fitnessCenter.information.name":               1,
		"fitnessCenter.information.image":              1,
		"_id":                                          1,
		"fitnessCenter.information.additionalText":     1,
		"fitnessCenter.review.totalReviews":            1,
		"fitnessCenter.review.avgRating":               1,
		"fitnessCenter.additionalServices.id":          1,
		"fitnessCenter.additionalServices.name":        1,
		"fitnessCenter.additionalServices.information": 1,
		"fitnessCenter.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
	}

	findOptions := options.FindOne().SetProjection(projection)

	var fitnessCenterData entity.ServiceEntity
	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&fitnessCenterData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fitnessCenter.GetFitnessCenterResDto{
			Status:  false,
			Message: "Failed to fetch fitnessCenter data: " + err.Error(),
		})
	}

	if fitnessCenterData.FitnessCenter == nil {
		return c.Status(fiber.StatusNotFound).JSON(fitnessCenter.GetFitnessCenterResDto{
			Status:  false,
			Message: "fitnessCenter data not found",
		})
	}

	servicesData := make([]fitnessCenter.AdditionalServices, 0)
	if fitnessCenterData.FitnessCenter != nil && len(fitnessCenterData.FitnessCenter.AdditionalServices) > 0 {
		for _, service := range fitnessCenterData.FitnessCenter.AdditionalServices {
			servicesData = append(servicesData, fitnessCenter.AdditionalServices{
				Id:          service.Id,
				Name:        service.Name,
				Information: service.Information,
			})
		}
	}

	fitnessCenterRes := fitnessCenter.GetFitnessCenterResDto{
		Status:  true,
		Message: "FitnessCenter data fetched successfully",
		Data: fitnessCenter.FitnessCenterResponse{
			Id:                     fitnessCenterData.Id,
			Image:                  fitnessCenterData.FitnessCenter.Information.Image,
			Name:                   fitnessCenterData.FitnessCenter.Information.Name,
			AboutUs:                fitnessCenterData.FitnessCenter.Information.AdditionalText,
			Address:                fitnessCenter.Address(fitnessCenterData.FitnessCenter.Information.Address),
			AdditionalServices:     servicesData,
			TotalReviews:           fitnessCenterData.FitnessCenter.Review.TotalReviews,
			AvgRating:              fitnessCenterData.FitnessCenter.Review.AvgRating,
			IsCustomerFamilyMember: isCustomerFamilyMember,
		},
	}

	return c.Status(fiber.StatusOK).JSON(fitnessCenterRes)
}
