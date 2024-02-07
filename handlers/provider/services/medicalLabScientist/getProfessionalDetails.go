package medicalLabScientist

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
// @Tags medicalLabScientist
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} services.GetMedicalLabScientistProfessionalDetailsResponseDto
// @Router /provider/services/get-medicalLabScientist-professional-details [get]
func FetchMedicalLabScientistProfessionalDetaiById(c *fiber.Ctx) error {

	var provider entity.ServiceEntity

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	serviceColl := database.GetCollection("service")

	projection := bson.M{
		"medicalLabScientist.professionalDetails.qualification":  1,
		"medicalLabScientist.professionalDetails.department":      1,
		"medicalLabScientist.professionalDetailsDocs.certificate": 1,
		"medicalLabScientist.professionalDetailsDocs.license":     1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	err := serviceColl.FindOne(ctx, bson.M{"_id": providerData.ProviderId}, findOptions).Decode(&provider)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.GetMedicalLabScientistProfessionalDetailsResponseDto{
				Status:  false,
				Message: "provider not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.GetMedicalLabScientistProfessionalDetailsResponseDto{
			Status:  false,
			Message: "Failed to fetch provider from MongoDB: " + err.Error(),
		})
	}

	var qualification string
	var department string
	var professionalLicense string
	var professionalCertificate string

	if provider.MedicalLabScientist != nil {
		qualification = provider.MedicalLabScientist.ProfessionalDetails.Qualification
		department = provider.MedicalLabScientist.ProfessionalDetails.Department
		professionalLicense = provider.MedicalLabScientist.ProfessionalDetailsDocs.License
		professionalCertificate = provider.MedicalLabScientist.ProfessionalDetailsDocs.Certificate
	}

	professionalDetailsRes := services.MedicalLabScientistProfessionalDetailsRes{
		Qualification:           qualification,
		ProfessionalLicense:     professionalLicense,
		ProfessionalCertificate: professionalCertificate,
		Department:              department,
	}

	return c.Status(fiber.StatusOK).JSON(services.GetMedicalLabScientistProfessionalDetailsResponseDto{
		Status:  true,
		Message: "Professional Details retrieved successfully",
		Data:    professionalDetailsRes,
	})
}
