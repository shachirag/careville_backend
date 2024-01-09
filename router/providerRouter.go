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
	provider.Put("/change-password/:id", jwt, providerAuthenticate.ChangeProviderPassword)
	provider.Get("/get-provider-info/:id", jwt, providerAuthenticate.FetchProviderById)
	provider.Put("/update-provider-info/:id", jwt, providerAuthenticate.UpdateProvider)

	// services
	provider.Post("/add-hospitalClinic", jwt, services.AddHospClinic)
	provider.Post("/add-laboratory", jwt, services.AddLaboratory)
	provider.Post("/add-fitness-center", jwt, services.AddFitnessCenter)
	provider.Post("/add-pharmacy", jwt, services.AddFitnessCenter)
	provider.Post("/add-doctor-profession", jwt, services.AddDoctorProfession)
	provider.Post("/add-nurse", jwt, services.AddNurse)
	provider.Post("/add-physiotherapist", jwt, services.AddPhysiotherapist)
	provider.Post("/add-medicalLab-scientist", jwt, services.AddMedicalLabScientist)
}
