package nurse

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// @Summary Fetch professionalDetails By ID
// @Description Fetch professionalDetails By ID
// @Tags nurse
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} services.GetNurseProfessionalDetailsResDto
// @Router /provider/services/get-nurse-professional-details [get]
func FetchProfessionalDetaiById(c *fiber.Ctx) error {

	var provider entity.ServiceEntity

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	serviceColl := database.GetCollection("service")

	projection := bson.M{
		"nurse.professionalDetails.qualifications":  1,
		"nurse.professionalDetailsDocs.certificate": 1,
		"nurse.professionalDetailsDocs.license":     1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	err := serviceColl.FindOne(ctx, bson.M{"_id": providerData.ProviderId}, findOptions).Decode(&provider)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.GetNurseProfessionalDetailsResDto{
				Status:  false,
				Message: "provider not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.GetNurseProfessionalDetailsResDto{
			Status:  false,
			Message: "Failed to fetch provider from MongoDB: " + err.Error(),
		})
	}

	var qualification string
	var professionalLicense string
	var professionalCertificate string

	if provider.Nurse != nil {
		qualification = provider.Nurse.ProfessionalDetails.Qualifications
		professionalLicense = provider.Nurse.ProfessionalDetailsDocs.License
		professionalCertificate = provider.Nurse.ProfessionalDetailsDocs.Certificate
	}

	professionalDetailsRes := services.PhysiotherapistDetailsRes{
		Qualification:           qualification,
		ProfessionalLicense:     professionalLicense,
		ProfessionalCertificate: professionalCertificate,
	}

	return c.Status(fiber.StatusOK).JSON(services.GetNurseProfessionalDetailsResDto{
		Status:  true,
		Message: "Professional Details retrieved successfully",
		Data:    professionalDetailsRes,
	})
}
