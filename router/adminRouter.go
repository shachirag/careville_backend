package router

import (
	"careville_backend/handlers"
	adminAuth "careville_backend/handlers/admin/adminAuthorization"
	dashboard "careville_backend/handlers/admin/dashboard"
	request "careville_backend/handlers/admin/requests"
	doctorProfession "careville_backend/handlers/admin/services/doctorProfession"
	fitnessCenter "careville_backend/handlers/admin/services/fitnessCenter"
	hospitals "careville_backend/handlers/admin/services/hospitals"
	laboratory "careville_backend/handlers/admin/services/laboratory"
	medicalLabScientist "careville_backend/handlers/admin/services/medicalLabScientist"
	nurse "careville_backend/handlers/admin/services/nurse"
	pharmacy "careville_backend/handlers/admin/services/pharmacy"
	physiotherapist "careville_backend/handlers/admin/services/physiotherapist"
	"careville_backend/middlewares"
	"os"

	"github.com/gofiber/fiber/v2"
)

func AdminSetupsRoutes(app *fiber.App) {

	app.Static("/", ".puplic")

	/* ---------- HEALTH ---------- */
	app.Get("/health", handlers.HealthCheck)

	/* ---------- Protected Routes ----- */
	secret := os.Getenv("JWT_SECRET_KEY")
	jwt := middlewares.NewAuthMiddleware(secret)

	// provider authentication
	admin := app.Group("/admin")
	admin.Post("/login", adminAuth.LoginAdmin)
	admin.Post("/verify-otp-for-login", adminAuth.VerifyOtpForLogin)
	admin.Post("/forgot-password", adminAuth.ForgotPassword)
	admin.Post("/verify-otp", adminAuth.VerifyOtp)
	admin.Put("/reset-password", adminAuth.ResetPasswordAfterOtp)
	admin.Post("/resend-otp", adminAuth.ResendOTP)

	adminProfile := admin.Group("/profile")
	adminProfile.Use(jwt, middlewares.AdminData)
	adminProfile.Get("/get-admin-info", adminAuth.FetchAdminById)
	adminProfile.Put("/change-password", adminAuth.ChangeAdminPassword)
	adminProfile.Put("/update-admin-info", adminAuth.UpdateAdmin)
	adminProfile.Put("/change-status/:id", adminAuth.ChangeStatus)

	// dashboard
	adminProfile.Get("/total-counts", jwt, dashboard.CountAll)

	requests := admin.Group("/requests")
	// requests.Use(jwt, middlewares.AdminData)
	requests.Get("/get-requests", request.FetchRequestsWithPagination)
	// services
	healthFacility := admin.Group("/healthFacility")
	healthFacility.Use(jwt, middlewares.AdminData)
	healthFacility.Get("/get-hospitals", hospitals.FetchHospitalsWithPagination)
	healthFacility.Get("/get-hospital/:id", hospitals.GetHospitalDetail)
	healthFacility.Get("/get-fitnessCenters", fitnessCenter.FetchFitnessCenterWithPagination)
	healthFacility.Get("/get-fitnessCenter/:id", fitnessCenter.GetFitnessCenterDetail)
	healthFacility.Get("/get-laboratories", laboratory.FetchLaboratoriesWithPagination)
	healthFacility.Get("/get-laboratory/:id", laboratory.GetLaboratoryDetail)
	healthFacility.Get("/get-pharmacies", pharmacy.FetchPharmacyWithPagination)
	healthFacility.Get("/get-pharmacy/:id", pharmacy.GetPharmacyDetail)

	healthProfessional := admin.Group("/healthProfessional")
	healthProfessional.Use(jwt, middlewares.AdminData)
	healthProfessional.Get("/get-doctors", doctorProfession.FetchDoctorsWithPagination)
	healthProfessional.Get("/get-doctor/:id", doctorProfession.GetDoctorDetail)
	healthProfessional.Get("/get-medicalLabScientists", medicalLabScientist.FetchMedicalLabScientistsWithPagination)
	healthProfessional.Get("/get-medicalLabScientist/:id", medicalLabScientist.GetMedicalLabScientistDetail)
	healthProfessional.Get("/get-nurses", nurse.FetchNurseWithPagination)
	healthProfessional.Get("/get-nurse/:id", nurse.GetNurseDetail)
	healthProfessional.Get("/get-physiotherapists", physiotherapist.FetchPhysiotherapistWithPagination)
	healthProfessional.Get("/get-physiotherapist/:id", physiotherapist.GetPhysiotherapistDetail)
}
