package router

import (
	"careville_backend/handlers"
	providerAuthenticate "careville_backend/handlers/provider/providerAuthentication"
	"careville_backend/handlers/provider/services"
	"careville_backend/middlewares"
	"os"

	"github.com/gofiber/fiber/v2"
)

func ProviderSetupsRoutes(app *fiber.App) {

	app.Static("/", ".puplic")

	/* ---------- HEALTH ---------- */
	app.Get("/health", handlers.HealthCheck)

	/* ---------- Protected Routes ----- */
	secret := os.Getenv("JWT_SECRET_KEY")
	jwt := middlewares.NewAuthMiddleware(secret)

	// provider authentication
	provider := app.Group("/provider")
	provider.Post("/signup", providerAuthenticate.SignupProvider)
	provider.Post("/verify-otp-for-signup", providerAuthenticate.VerifyOtpForSignup)
	provider.Post("/login", providerAuthenticate.LoginProvider)
	provider.Post("/forgot-password", providerAuthenticate.ForgotPassword)
	provider.Post("/verify-otp", providerAuthenticate.VerifyOtp)
	provider.Put("/reset-password", providerAuthenticate.ResetPasswordAfterOtp)
	providerProfile := provider.Group("/profile")
	providerProfile.Use(jwt, middlewares.ProviderData)
	providerProfile.Put("/change-password", providerAuthenticate.ChangeProviderPassword)
	providerProfile.Get("/get-provider-info", providerAuthenticate.FetchProviderById)
	providerProfile.Put("/update-provider-info", providerAuthenticate.UpdateProvider)
	providerProfile.Put("/update-profile-image", providerAuthenticate.UpdateImage)

	providerServices := provider.Group("/services")
	providerServices.Use(jwt, middlewares.ProviderData)
	providerServices.Post("/add-hospitalClinic", services.AddHospClinic)
	providerServices.Post("/add-laboratory", services.AddLaboratory)
	providerServices.Post("/add-fitness-center", services.AddFitnessCenter)
	providerServices.Post("/add-pharmacy", services.AddPharmacy)
	providerServices.Post("/add-doctor-profession", services.AddDoctorProfession)
	providerServices.Post("/add-nurse", services.AddNurse)
	providerServices.Post("/add-physiotherapist", services.AddPhysiotherapist)
	providerServices.Post("/add-medicalLab-scientist", services.AddMedicalLabScientist)
	providerServices.Get("/get-misc-data", services.FetchAllMiscData)
	// provider.Put("/change-status/:id", jwt, services.ChangeStatus)
	providerServices.Get("/get-status", services.FetchStatusById)
	providerServices.Get("/get-all-doctors", services.GetAllDoctors)
	providerServices.Get("/get-investigations", services.GetInvestigations)
	providerServices.Post("/add-more-doctor", services.AddMoreDoctors)
	providerServices.Get("/get-doctor-info/:doctorId", services.GetDoctorsInfo)
	providerServices.Get("/get-investigation-info/:investigationId", services.GetInvesitagtionInfo)
	providerServices.Put("/update-doctor-image/:doctorId", services.UpdateDoctorImage)
	providerServices.Post("/add-other-services", services.AddServices)
	providerServices.Put("/update-doctor-info/:doctorId", services.UpdateDoctorInfo)

}
