package laboratory

import (
	"careville_backend/database"
	laboratory "careville_backend/dto/admin/services/laboratories"
	"careville_backend/entity"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary get-laboratory detail
// @Description get-laboratory detail
// @Tags admin laboratory
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} laboratory.GetLaboratoryDetailResDto
// @Router /admin/healthFacility/get-laboratory/{id} [get]
func GetLaboratoryDetail(c *fiber.Ctx) error {
	ctx := context.TODO()

	var service entity.ServiceEntity

	serviceColl := database.GetCollection("service")

	idParam := c.Params("id")
	hospitalID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(laboratory.GetLaboratoryDetailResDto{
			Status:  false,
			Message: "Invalid appointment ID",
		})
	}

	filter := bson.M{
		"_id": hospitalID,
	}

	projection := bson.M{
		"_id":                                   1,
		"facilityOrProfession":                  1,
		"role":                                  1,
		"profileId":                             1,
		"serviceStatus":                         1,
		"user.firstName":                        1,
		"user.lastName":                         1,
		"user.email":                            1,
		"user.phoneNumber.dialCode":             1,
		"user.phoneNumber.countryCode":          1,
		"user.phoneNumber.number":               1,
		"laboratory.documents.certificate":      1,
		"laboratory.documents.license":          1,
		"laboratory.information.name":           1,
		"laboratory.information.image":          1,
		"laboratory.information.additionalText": 1,
		"laboratory.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
		"laboratory.investigations.id":          1,
		"laboratory.investigations.name":        1,
		"laboratory.investigations.type":        1,
		"laboratory.investigations.information": 1,
		"laboratory.investigations.price":       1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(laboratory.GetLaboratoryDetailResDto{
				Status:  false,
				Message: "Other service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(laboratory.GetLaboratoryDetailResDto{
			Status:  false,
			Message: "Failed to fetch other service from MongoDB: " + err.Error(),
		})
	}

	if service.Laboratory == nil {
		return c.Status(fiber.StatusOK).JSON(laboratory.GetLaboratoryDetailResDto{
			Status:  false,
			Message: "Laboratory information not found.",
		})
	}

	investigationData := make([]laboratory.Investigations, 0)
	if service.Laboratory != nil && len(service.Laboratory.Investigations) > 0 {
		for _, investigation := range service.Laboratory.Investigations {
			investigationData = append(investigationData, laboratory.Investigations{
				Id:          investigation.Id,
				Type:        investigation.Type,
				Name:        investigation.Name,
				Information: investigation.Information,
				Price:       investigation.Price,
			})
		}
	}

	var laboratoryImage string
	var laboratoryName string
	var additionalText string
	var license string
	var certificate string
	var laboratoryAddress laboratory.Address
	if service.Laboratory != nil {
		laboratoryName = service.Laboratory.Information.Name
		laboratoryImage = service.Laboratory.Information.Image
		additionalText = service.Laboratory.Information.AdditionalText
		laboratoryAddress = laboratory.Address(service.Laboratory.Information.Address)
		license = service.Laboratory.Documents.License
		certificate = service.Laboratory.Documents.Certificate
	}

	response := laboratory.GetLaboratoryDetailResDto{
		Status:  true,
		Message: "Data fetched retrieved successfully",
		Data: laboratory.GetLaboratoryDetailRes{
			Id:                   service.Id,
			FacilityOrProfession: service.FacilityOrProfession,
			Role:                 service.Role,
			ProfileId:            service.ProfileId,
			ServiceStatus:        service.ServiceStatus,
			User: laboratory.User{
				FirstName: service.User.FirstName,
				LastName:  service.User.LastName,
				Email:     service.User.Email,
				PhoneNumber: laboratory.PhoneNumber{
					DialCode: service.User.PhoneNumber.DialCode,
					Number:   service.User.PhoneNumber.Number,
				},
			},
			Investigations: investigationData,
			LaboratoryInformation: laboratory.LaboratoryInformation{
				Name:           laboratoryName,
				Image:          laboratoryImage,
				AdditionalText: additionalText,
				Address:        laboratoryAddress,
			},
			Documents: laboratory.Documents{
				License:     license,
				Certificate: certificate,
			},
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
