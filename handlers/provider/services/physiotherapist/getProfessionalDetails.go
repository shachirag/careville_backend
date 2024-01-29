package physiotherapist

// import (
// 	"careville_backend/database"
// 	providerMiddleware "careville_backend/dto/provider/middleware"
// 	providerAuth "careville_backend/dto/provider/providerAuth"
// 	"careville_backend/entity"

// 	"github.com/gofiber/fiber/v2"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// // @Summary Fetch professionalDetails By ID
// // @Description Fetch professionalDetails By ID
// // @Tags physiotherapist
// // @Accept application/json
// //
// //	@Param Authorization header	string true	"Authentication header"
// //
// // @Produce json
// // @Success 200 {object} providerAuth.GetProviderResDto
// // @Router /provider/services/get-physiotherapist-professional-details [get]
// func FetchProfessionalDetaiById(c *fiber.Ctx) error {

// 	var provider entity.ServiceEntity

// 	// Get provider data from middleware
// 	providerData := providerMiddleware.GetProviderMiddlewareData(c)

// 	serviceColl := database.GetCollection("service")

// 	err := serviceColl.FindOne(ctx, bson.M{"_id": providerData.ProviderId}).Decode(&provider)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return c.Status(fiber.StatusNotFound).JSON(providerAuth.GetProviderResDto{
// 				Status:  false,
// 				Message: "provider not found",
// 			})
// 		}
// 		return c.Status(fiber.StatusInternalServerError).JSON(providerAuth.GetProviderResDto{
// 			Status:  false,
// 			Message: "Failed to fetch provider from MongoDB: " + err.Error(),
// 		})
// 	}

// 	var qualification string
// 	var speciality string
// 	var serviceFees bool
// 	var professionalLicense string
// 	var professionalCertificate string

// 	if provider.Role == "healthProfessional" && provider.FacilityOrProfession == "physiotherapist" {
// 		qualification = provider.Physiotherapist.ProfessionalDetails.Qualifications
// 		serviceFees = provider.Physiotherapist.ServiceAndSchedule.
// 		name = provider.Physiotherapist.Information.Name
// 	}

// 	providerRes := providerAuth.ProviderResDto{
// 		User: providerAuth.UserData{
// 			Role: providerAuth.Role{
// 				Role:                 provider.Role,
// 				FacilityOrProfession: provider.FacilityOrProfession,
// 				ServiceStatus:        provider.ServiceStatus,
// 				Image:                image,
// 				Name:                 name,
// 				IsEmergencyAvailable: isEmergencyAvailable,
// 			},
// 			User: providerAuth.User{
// 				Id:          provider.Id,
// 				FirstName:   provider.User.FirstName,
// 				LastName:    provider.User.LastName,
// 				Email:       provider.User.Email,
// 				PhoneNumber: providerAuth.PhoneNumber(provider.User.PhoneNumber),
// 				Notification: providerAuth.Notification{
// 					DeviceToken: provider.User.Notification.DeviceToken,
// 					DeviceType:  provider.User.Notification.DeviceType,
// 					IsEnabled:   provider.User.Notification.IsEnabled,
// 				},
// 				CreatedAt: provider.CreatedAt,
// 				UpdatedAt: provider.UpdatedAt,
// 			},
// 		},
// 		AdditionalInformation: providerAuth.AdditionalInformation{
// 			AdditionalDetails: additionalDetails,
// 			Address:           address,
// 			Documents: providerAuth.Documents{
// 				Certificate: certificate,
// 				License:     license,
// 			},
// 		},
// 	}

// 	return c.Status(fiber.StatusOK).JSON(providerAuth.GetProviderResDto{
// 		Status:   true,
// 		Message:  "provider data retrieved successfully",
// 		Provider: providerRes,
// 	})
// }
