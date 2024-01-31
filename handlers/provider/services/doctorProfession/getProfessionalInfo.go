package doctorProfession

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
// @Tags doctorProfession
// @Accept application/json
//
//	@Param Authorization header	string true	"Authentication header"
//
// @Produce json
// @Success 200 {object} services.GetDoctorProfessionProfessionalDetailsResDto
// @Router /provider/services/get-doctorProfession-professional-details [get]
func FetchDoctorProfessionProfessionalDetaiById(c *fiber.Ctx) error {

	var provider entity.ServiceEntity

	// Get provider data from middleware
	providerData := providerMiddleware.GetProviderMiddlewareData(c)

	serviceColl := database.GetCollection("service")

	projection := bson.M{
		"doctor.additionalServices.qualifications":   1,
		"doctor.additionalServices.speciality":       1,
		"doctor.professionalDetailsDocs.certificate": 1,
		"doctor.professionalDetailsDocs.license":     1,
		"doctor.schedule.consultationFees":           1,
	}

	findOptions := options.FindOne().SetProjection(projection)

	err := serviceColl.FindOne(ctx, bson.M{"_id": providerData.ProviderId}, findOptions).Decode(&provider)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(services.GetDoctorProfessionProfessionalDetailsResDto{
				Status:  false,
				Message: "provider not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(services.GetDoctorProfessionProfessionalDetailsResDto{
			Status:  false,
			Message: "Failed to fetch provider from MongoDB: " + err.Error(),
		})
	}

	var qualification string
	var speciality string
	var professionalLicense string
	var professionalCertificate string
	var consultingFees float64

	if provider.Doctor != nil {
		qualification = provider.Doctor.AdditionalServices.Qualifications
		speciality = provider.Doctor.AdditionalServices.Speciality
		professionalLicense = provider.Doctor.ProfessionalDetailsDocs.License
		professionalCertificate = provider.Doctor.ProfessionalDetailsDocs.Certificate
		consultingFees = provider.Doctor.Schedule.ConsultationFees
	}

	professionalDetailsRes := services.DoctorProfessionProfessionDetailsRes{
		Qualification:           qualification,
		ProfessionalLicense:     professionalLicense,
		ProfessionalCertificate: professionalCertificate,
		ConsultingFees:          consultingFees,
		Speciality:              speciality,
	}

	return c.Status(fiber.StatusOK).JSON(services.GetDoctorProfessionProfessionalDetailsResDto{
		Status:  true,
		Message: "Professional Details retrieved successfully",
		Data:    professionalDetailsRes,
	})
}
