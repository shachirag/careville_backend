package router

import (
	"careville_backend/handlers"
	providerAuthenticate "careville_backend/handlers/provider/providerAuthentication"
	common "careville_backend/handlers/provider/services/commonApis"
	"careville_backend/handlers/provider/services/doctorProfession"
	"careville_backend/handlers/provider/services/fitnessCenter"
	"careville_backend/handlers/provider/services/hospClinic"
	"careville_backend/handlers/provider/services/laboratory"
	"careville_backend/handlers/provider/services/medicalLabScientist"
	"careville_backend/handlers/provider/services/nurse"
	"careville_backend/handlers/provider/services/pharmacy"
	"careville_backend/handlers/provider/services/physiotherapist"
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

	providerServices.Post("/add-laboratory", laboratory.AddLaboratory)
	providerServices.Post("/add-hospitalClinic", hospClinic.AddHospClinic)
	providerServices.Post("/add-fitness-center", fitnessCenter.AddFitnessCenter)
	providerServices.Post("/add-pharmacy", pharmacy.AddPharmacy)
	providerServices.Post("/add-doctor-profession", doctorProfession.AddDoctorProfession)
	providerServices.Post("/add-nurse", nurse.AddNurse)
	providerServices.Post("/add-physiotherapist", physiotherapist.AddPhysiotherapist)
	providerServices.Post("/add-medicalLab-scientist", medicalLabScientist.AddMedicalLabScientist)

	providerServices.Get("/get-misc-data", common.FetchAllMiscData)
	// provider.Put("/change-status/:id", jwt, services.ChangeStatus)
	providerServices.Get("/get-status", common.FetchStatusById)

	providerServices.Get("/get-all-doctors", hospClinic.GetAllDoctors)
	providerServices.Get("/get-all-trainers", fitnessCenter.GetAllTrainers)
	providerServices.Get("/get-pharmacy-other-services", pharmacy.GetOtherServices)
	providerServices.Get("/get-fitness-other-services", fitnessCenter.GetOtherServices)
	providerServices.Get("/get-all-subscriptions", fitnessCenter.GetAllSubscriptions)
	providerServices.Get("/get-hospClinic-other-services", hospClinic.GetOtherServices)
	providerServices.Get("/get-doctor-info/:doctorId", hospClinic.GetDoctorsInfo)
	providerServices.Get("/get-pharmacy-other-service-info/:otherServiceId", pharmacy.GetOtherServiceInfo)
	providerServices.Get("/get-trainer-info/:trainerId", fitnessCenter.GetTrainerInfo)
	providerServices.Get("/get-fitness-other-service-info/:otherServiceId", fitnessCenter.GetOtherServiceInfo)
	providerServices.Put("/update-pharmacy-other-service-info/:otherServiceId", pharmacy.UpdateOtherServiceInfo)
	providerServices.Post("/add-pharmacy-other-service", pharmacy.AddOtherServices)
	// providerServices.Put("/update-doctor-image/:doctorId", hospClinic.UpdateDoctorImage)
	providerServices.Post("/add-hospClinic-other-services", hospClinic.AddServices)
	providerServices.Put("/update-doctor-info/:doctorId", hospClinic.UpdateDoctorInfo)
	providerServices.Put("/update-trainer-info/:trainerId", fitnessCenter.UpdateTrainerInfo)
	providerServices.Put("/update-fitnessCenter-other-service-info/:otherServiceId", fitnessCenter.UpdateOtherServiceInfo)
	providerServices.Get("/get-subscription-info/:subscriptionId", fitnessCenter.GetSubscriptionInfo)
	providerServices.Post("/add-more-doctor", hospClinic.AddMoreDoctors)
	providerServices.Get("/get-investigations", laboratory.GetInvestigations)

	providerServices.Put("/change-notification", common.ProviderNotification)
	providerServices.Put("/currently-available", common.ProviderCurrentlyAvailable)
	providerServices.Post("/add-more-investigation", laboratory.AddMoreInvestigstions)
	providerServices.Post("/add-more-trainer", fitnessCenter.AddMoreTrainers)
	providerServices.Post("/add-more-subscription", fitnessCenter.AddMoreSubscriptions)
	providerServices.Post("/add-fitnessCenter-other-service", fitnessCenter.AddOtherServices)
	providerServices.Get("/get-investigation-info/:investigationId", laboratory.GetInvesitagtionInfo)
	providerServices.Put("/update-investigation-info/:investigationId", laboratory.UpdateinvestigationInfo)

}
