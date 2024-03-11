package pharmacy

import (
	"careville_backend/database"
	pharmacy "careville_backend/dto/admin/services/pharmacy"
	"careville_backend/entity"
	"context"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary get-pharmacy detail
// @Description get-pharmacy detail
// @Tags admin pharmacy
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} pharmacy.GetPharmacyDetailResDto
// @Router /admin/healthFacility/get-pharmacy/{id} [get]
func GetPharmacyDetail(c *fiber.Ctx) error {
	ctx := context.TODO()

	var service entity.ServiceEntity

	serviceColl := database.GetCollection("service")

	idParam := c.Params("id")
	hospitalID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(pharmacy.GetPharmacyDetailResDto{
			Status:  false,
			Message: "Invalid appointment ID",
		})
	}

	filter := bson.M{
		"_id": hospitalID,
	}

	projection := bson.M{
		"_id":                                 1,
		"facilityOrProfession":                1,
		"role":                                1,
		"profileId":                           1,
		"serviceStatus":                       1,
		"user.firstName":                      1,
		"user.lastName":                       1,
		"user.email":                          1,
		"user.phoneNumber.dialCode":           1,
		"user.phoneNumber.countryCode":        1,
		"user.phoneNumber.number":             1,
		"pharmacy.documents.certificate":      1,
		"pharmacy.documents.license":          1,
		"pharmacy.information.name":           1,
		"pharmacy.information.image":          1,
		"pharmacy.information.additionalText": 1,
		"pharmacy.information.address": bson.M{
			"coordinates": 1,
			"type":        1,
			"add":         1,
		},
		"pharmacy.additionalServices.id":          1,
		"pharmacy.additionalServices.name":        1,
		"pharmacy.additionalServices.information": 1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	err = serviceColl.FindOne(ctx, filter, findOptions).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(pharmacy.GetPharmacyDetailResDto{
				Status:  false,
				Message: "Other service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(pharmacy.GetPharmacyDetailResDto{
			Status:  false,
			Message: "Failed to fetch other service from MongoDB: " + err.Error(),
		})
	}

	if service.Pharmacy == nil {
		return c.Status(fiber.StatusOK).JSON(pharmacy.GetPharmacyDetailResDto{
			Status:  false,
			Message: "Pharmacy information not found.",
		})
	}

	otherServicesData := make([]pharmacy.OtherServices, 0)
	if service.Pharmacy != nil && len(service.Pharmacy.AdditionalServices) > 0 {
		for _, service := range service.Pharmacy.AdditionalServices {
			otherServicesData = append(otherServicesData, pharmacy.OtherServices{
				Id:          service.Id,
				Name:        service.Name,
				Information: service.Information,
			})
		}
	}

	var pharmacyImage string
	var pharmacyName string
	var additionalText string
	var license string
	var certificate string
	var pharmacyAddress pharmacy.Address
	if service.Pharmacy != nil {
		pharmacyName = service.Pharmacy.Information.Name
		pharmacyImage = service.Pharmacy.Information.Image
		additionalText = service.Pharmacy.Information.AdditionalText
		pharmacyAddress = pharmacy.Address(service.Pharmacy.Information.Address)
		license = service.Pharmacy.Documents.License
		certificate = service.Pharmacy.Documents.Certificate
	}

	response := pharmacy.GetPharmacyDetailResDto{
		Status:  true,
		Message: "Data fetched retrieved successfully",
		Data: pharmacy.GetPharmacyDetailRes{
			Id:                   service.Id,
			FacilityOrProfession: service.FacilityOrProfession,
			Role:                 service.Role,
			ProfileId:            service.ProfileId,
			ServiceStatus:        service.ServiceStatus,
			User: pharmacy.User{
				FirstName: service.User.FirstName,
				LastName:  service.User.LastName,
				Email:     service.User.Email,
				PhoneNumber: pharmacy.PhoneNumber{
					DialCode: service.User.PhoneNumber.DialCode,
					Number:   service.User.PhoneNumber.Number,
				},
			},
			OtherServices: otherServicesData,
			PharmacyInformation: pharmacy.PharmacyInformation{
				Name:           pharmacyName,
				Image:          pharmacyImage,
				AdditionalText: additionalText,
				Address:        pharmacyAddress,
			},
			Documents: pharmacy.Documents{
				License:     license,
				Certificate: certificate,
			},
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
