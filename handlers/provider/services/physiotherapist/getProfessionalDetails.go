package physiotherapist

import (
	"careville_backend/database"
	providerMiddleware "careville_backend/dto/provider/middleware"
	"careville_backend/dto/provider/services"
	"careville_backend/entity"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// @Summary Fetch professionalDetails By ID
// @Description Fetch professionalDetails By ID
// @Tags physiotherapist
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} services.GetPhysiotherapistProfessionalDetailsResDto
// @Router /provider/services/get-physiotherapist-professional-details [get]
func FetchProfessionalDetaiById(c *fiber.Ctx) error {

	var provider entity.ServiceEntity

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	serviceColl := database.GetCollection("service")

	err := serviceColl.FindOne(ctx, bson.M{"_id": providerData.ProviderId}).Decode(&provider)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.GetPhysiotherapistProfessionalDetailsResDto{
				Status:  false,
				Message: "provider not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.GetPhysiotherapistProfessionalDetailsResDto{
			Status:  false,
			Message: "Failed to fetch provider from MongoDB: " + err.Error(),
		})
	}

	var qualification string
	var professionalLicense string
	var professionalCertificate string

	qualification = provider.Physiotherapist.ProfessionalDetails.Qualifications
	professionalLicense = provider.Physiotherapist.ProfessionalDetailsDocs.License
	professionalCertificate = provider.Physiotherapist.ProfessionalDetailsDocs.Certificate

	professionalDetailsRes := services.PhysiotherapistDetailsRes{
		Qualification:           qualification,
		ProfessionalLicense:     professionalLicense,
		ProfessionalCertificate: professionalCertificate,
	}

	return c.Status(fiber.StatusOK).JSON(services.GetPhysiotherapistProfessionalDetailsResDto{
		Status:  true,
		Message: "Professional Details retrieved successfully",
		Data:    professionalDetailsRes,
	})
}
