package fitnessCenter

import (
	"careville_backend/database"
	fitnessCenter "careville_backend/dto/admin/services/fitnessCenter"
	"careville_backend/entity"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary get-fitnessCenter detail
// @Description get-fitnessCenter detail
// @Tags admin fitnessCenter
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} fitnessCenter.GetFitnessCenterDetailResDto
// @Router /admin/healthFacility/get-fitnessCenter/{id} [get]
func GetFitnessCenterDetail(c *fiber.Ctx) error {
	ctx := context.TODO()

	var service entity.ServiceEntity

	serviceColl := database.GetCollection("service")

	idParam := c.Params("id")
	hospitalID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fitnessCenter.GetFitnessCenterDetailResDto{
			Status:  false,
			Message: "Invalid appointment ID",
		})
	}

	filter := bson.M{
		"_id": hospitalID,
	}

	projection := bson.M{
		"_id":                                      1,
		"facilityOrProfession":                     1,
		"role":                                     1,
		"profileId":                                1,
		"serviceStatus":                            1,
		"user.firstName":                           1,
		"user.lastName":                            1,
		"user.email":                               1,
		"user.phoneNumber.dialCode":                1,
		"user.phoneNumber.countryCode":             1,
		"user.phoneNumber.number":                  1,
		"fitnessCenter.documents.certificate":      1,
		"fitnessCenter.documents.license":          1,
		"fitnessCenter.information.name":           1,
		"fitnessCenter.information.image":          1,
		"fitnessCenter.information.additionalText": 1,
		"fitnessCenter.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
		"fitnessCenter.trainers.id":                    1,
		"fitnessCenter.trainers.name":                  1,
		"fitnessCenter.trainers.category":              1,
		"fitnessCenter.trainers.information":           1,
		"fitnessCenter.trainers.price":                 1,
		"fitnessCenter.additionalServices.id":          1,
		"fitnessCenter.additionalServices.information": 1,
		"fitnessCenter.additionalServices.name":        1,
		"fitnessCenter.subscription.id":                1,
		"fitnessCenter.subscription.type":              1,
		"fitnessCenter.subscription.details":           1,
		"fitnessCenter.subscription.price":             1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fitnessCenter.GetFitnessCenterDetailResDto{
				Status:  false,
				Message: "Other service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fitnessCenter.GetFitnessCenterDetailResDto{
			Status:  false,
			Message: "Failed to fetch other service from MongoDB: " + err.Error(),
		})
	}

	if service.FitnessCenter == nil {
		return c.Status(fiber.StatusOK).JSON(fitnessCenter.GetFitnessCenterDetailResDto{
			Status:  false,
			Message: "FitnessCenter information not found.",
		})
	}

	trainersData := make([]fitnessCenter.Trainers, 0)
	if service.FitnessCenter != nil && len(service.FitnessCenter.Trainers) > 0 {
		for _, trainer := range service.FitnessCenter.Trainers {
			trainersData = append(trainersData, fitnessCenter.Trainers{
				Id:          trainer.Id,
				Category:    trainer.Category,
				Name:        trainer.Name,
				Information: trainer.Information,
				Price:       trainer.Price,
			})
		}
	}

	subscriptionData := make([]fitnessCenter.Subscripions, 0)
	if service.FitnessCenter != nil && len(service.FitnessCenter.Subscription) > 0 {
		for _, subscription := range service.FitnessCenter.Subscription {
			subscriptionData = append(subscriptionData, fitnessCenter.Subscripions{
				Id:      subscription.Id,
				Type:    subscription.Type,
				Details: subscription.Details,
				Price:   subscription.Price,
			})
		}
	}

	otherServicesData := make([]fitnessCenter.OtherServices, 0)
	if service.FitnessCenter != nil && len(service.FitnessCenter.AdditionalServices) > 0 {
		for _, service := range service.FitnessCenter.AdditionalServices {
			otherServicesData = append(otherServicesData, fitnessCenter.OtherServices{
				Id:          service.Id,
				Information: service.Information,
				Name:        service.Name,
			})
		}
	}

	var fitnessCenterImage string
	var fitnessCenterName string
	var additionalText string
	var license string
	var certificate string
	var fitnessCenterAddress fitnessCenter.Address
	if service.FitnessCenter != nil {
		fitnessCenterName = service.FitnessCenter.Information.Name
		fitnessCenterImage = service.FitnessCenter.Information.Image
		additionalText = service.FitnessCenter.Information.AdditionalText
		fitnessCenterAddress = fitnessCenter.Address(service.FitnessCenter.Information.Address)
		license = service.FitnessCenter.Documents.License
		certificate = service.FitnessCenter.Documents.Certificate
	}

	response := fitnessCenter.GetFitnessCenterDetailResDto{
		Status:  true,
		Message: "Data fetched retrieved successfully",
		Data: fitnessCenter.GetFitnessCenterDetailRes{
			Id:                   service.Id,
			FacilityOrProfession: service.FacilityOrProfession,
			Role:                 service.Role,
			ProfileId:            service.ProfileId,
			ServiceStatus:        service.ServiceStatus,
			User: fitnessCenter.User{
				FirstName: service.User.FirstName,
				LastName:  service.User.LastName,
				Email:     service.User.Email,
				PhoneNumber: fitnessCenter.PhoneNumber{
					DialCode: service.User.PhoneNumber.DialCode,
					Number:   service.User.PhoneNumber.Number,
				},
			},
			Trainers:      trainersData,
			Subscripions:  subscriptionData,
			OtherServices: otherServicesData,
			FitnessCenterInformation: fitnessCenter.FitnessCenterInformation{
				Name:           fitnessCenterName,
				Image:          fitnessCenterImage,
				AdditionalText: additionalText,
				Address:        fitnessCenterAddress,
			},
			Documents: fitnessCenter.Documents{
				License:     license,
				Certificate: certificate,
			},
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
