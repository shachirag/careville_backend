package router

import (
	"careville_backend/handlers"
	customerAuth "careville_backend/handlers/customer/customerAuthentication"
	"careville_backend/handlers/customer/doctorProfession"
	"careville_backend/handlers/customer/fitnessCenter"
	hospitals "careville_backend/handlers/customer/hospitals"
	laboratory "careville_backend/handlers/customer/laboratory"
	"careville_backend/handlers/customer/medicalLabScientist"
	"careville_backend/handlers/customer/nurse"
	"careville_backend/handlers/customer/pharmacy"
	"careville_backend/handlers/customer/physiotherapist"
	"careville_backend/middlewares"
	"os"

	"github.com/gofiber/fiber/v2"
)

func CustomerSetupsRoutes(app *fiber.App) {

	app.Static("/", ".puplic")

	/* ---------- HEALTH ---------- */
	app.Get("/health", handlers.HealthCheck)

	/* ---------- Protected Routes ----- */
	secret := os.Getenv("JWT_SECRET_KEY")
	jwt := middlewares.NewAuthMiddleware(secret)

	customer := app.Group("/customer")
	customer.Post("/login", customerAuth.LoginCustomer)
	customer.Post("/forgot-password", customerAuth.ForgotPassword)
	customer.Post("/verify-otp", customerAuth.VerifyOtp)
	customer.Put("/reset-password", customerAuth.ResetPasswordAfterOtp)
	customer.Post("/signup", customerAuth.SignupCustomer)
	customer.Post("/verify-otp-for-signup", customerAuth.VerifyOtpForSignup)

	customerProfile := customer.Group("/profile")
	customerProfile.Use(jwt, middlewares.CustomerData)

	customerProfile.Get("/get-customer-info/:id", customerAuth.GetCustomer)
	customerProfile.Put("/change-password", customerAuth.ChangeCustomerPassword)
	customerProfile.Put("/change-notification", customerAuth.CustomerNotification)
	customerProfile.Post("/add-more-family-member", customerAuth.AddMoreMembers)
	customerProfile.Get("/get-members", customerAuth.GetMembers)
	customerProfile.Put("/edit-customer-info", customerAuth.UpdateCustomer)

	healthFacility := customer.Group("/healthFacility")
	healthFacility.Use(jwt, middlewares.CustomerData)
	healthFacility.Get("/get-hospitals", hospitals.FetchHospitalsWithPagination)
	healthFacility.Get("/get-hospital/:id", hospitals.GetHospitalByID)
	healthFacility.Get("/get-laboratories", laboratory.FetchLaboratoriesWithPagination)
	healthFacility.Get("/get-laboratory/:id", laboratory.GetLaboratoryByID)
	healthFacility.Get("/get-fitnessCenters", fitnessCenter.FetchFitnessCenterWithPagination)
	healthFacility.Get("/get-fitnessCenter/:id", fitnessCenter.GetFitnessCenterByID)
	healthFacility.Get("/get-pharmacies", pharmacy.FetchPharmacyWithPagination)
	healthFacility.Get("/get-pharmacy/:id", pharmacy.GetPharmacyByID)

	healthProfessional := customer.Group("/healthProfessional")
	healthProfessional.Use(jwt, middlewares.CustomerData)

	healthProfessional.Get("/get-physiotherapists", physiotherapist.FetchPhysiotherapistWithPagination)
	healthProfessional.Get("/get-physiotherapist/:id", physiotherapist.GetPhysiotherapistByID)
	healthProfessional.Get("/get-medicalLabScientists", medicalLabScientist.FetchMedicalLabScientistsWithPagination)
	healthProfessional.Get("/get-medicalLabScientist/:id", medicalLabScientist.GetMedicalLabScientistByID)
	healthProfessional.Get("/get-nurses", nurse.FetchNurseWithPagination)
	healthProfessional.Get("/get-nurse/:id", nurse.GetNurseByID)
	healthProfessional.Get("/get-doctors", doctorProfession.FetchDoctorsWithPagination)
	healthProfessional.Get("/get-doctor/:id", doctorProfession.GetDoctorProfessionByID)
}
