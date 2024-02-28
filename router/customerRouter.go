package router

import (
	"careville_backend/handlers"
	common "careville_backend/handlers/customer/commonApis"
	customerAuth "careville_backend/handlers/customer/customerAuthentication"
	"careville_backend/handlers/customer/doctorProfession"
	emergency "careville_backend/handlers/customer/emergency"
	"careville_backend/handlers/customer/fitnessCenter"
	hospitals "careville_backend/handlers/customer/hospitals"
	laboratory "careville_backend/handlers/customer/laboratory"
	"careville_backend/handlers/customer/medicalLabScientist"
	"careville_backend/handlers/customer/nurse"
	"careville_backend/handlers/customer/pharmacy"
	"careville_backend/handlers/customer/physiotherapist"

	// "careville_backend/handlers/customer/reviews"
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
	// customer.Post("/add-review", jwt, reviews.AddReview)

	customerProfile := customer.Group("/profile")
	customerProfile.Use(jwt, middlewares.CustomerData)

	customerProfile.Get("/get-customer-info", customerAuth.GetCustomer)
	customerProfile.Put("/change-password", customerAuth.ChangeCustomerPassword)
	customerProfile.Put("/change-notification", customerAuth.CustomerNotification)
	customerProfile.Post("/add-more-family-member", customerAuth.AddMoreMembers)
	customerProfile.Get("/get-members", customerAuth.GetMembers)
	customerProfile.Put("/edit-customer-info", customerAuth.UpdateCustomer)

	healthFacility := customer.Group("/healthFacility")
	healthFacility.Use(jwt, middlewares.CustomerData)
	healthFacility.Get("/get-health-facilities", common.GetHealthFacilties)
	healthFacility.Post("/add-hospClinic-appointment", hospitals.AddHospClinicAppointment)
	healthFacility.Get("/get-hospitals", hospitals.GetHospitals)
	healthFacility.Get("/get-hospital/:id", hospitals.GetHospitalByID)
	healthFacility.Get("/get-all-doctors", hospitals.GetAllDoctors)
	healthFacility.Get("/get-all-available-times", hospitals.GetAvailableSlots)
	healthFacility.Post("/add-laboratory-appointment", laboratory.AddLaboratoryAppointment)
	healthFacility.Get("/get-laboratories", laboratory.Getlaboratory)
	healthFacility.Get("/get-investigations", laboratory.GetInvestigations)
	healthFacility.Get("/get-laboratory/:id", laboratory.GetLaboratoryByID)
	healthFacility.Post("/add-fitnessCenter-appointment", fitnessCenter.AddFitnessCenterAppointment)
	healthFacility.Get("/get-fitnessCenters", fitnessCenter.GetFitnessCenter)
	healthFacility.Get("/get-subscriptions/:id", fitnessCenter.GetAllSubscriptions)
	healthFacility.Get("/get-fitnessCenter/:id", fitnessCenter.GetFitnessCenterByID)
	healthFacility.Get("/get-all-trainers", fitnessCenter.GetAllTrainers)
	healthFacility.Post("/add-pharmacy-drug", pharmacy.AddPharmacyDrugs)
	healthFacility.Get("/get-pharmacies", pharmacy.GetPharmacy)
	healthFacility.Get("/get-pharmacy/:id", pharmacy.GetPharmacyByID)

	healthProfessional := customer.Group("/healthProfessional")
	healthProfessional.Use(jwt, middlewares.CustomerData)
	healthProfessional.Get("/get-health-professionals", common.GetHealthProfessionals)
	healthProfessional.Post("/add-physiotherapist-appointment", physiotherapist.AddPhysiotherapistAppointment)
	healthProfessional.Get("/get-physiotherapists", physiotherapist.GetPhysiotherapists)
	healthProfessional.Get("/get-physiotherapist/:id", physiotherapist.GetPhysiotherapistByID)
	healthProfessional.Get("/get-physiotherapist-services", physiotherapist.GetPhysiotherapistServices)
	healthProfessional.Get("/get-physiotherapist-available-slots", physiotherapist.GetPhysiotherapistAvailableTimes)
	healthProfessional.Post("/add-medicalLabScientist-appointment", medicalLabScientist.AddMedicalLabScientistAppointment)
	healthProfessional.Get("/get-medicalLabScientists", medicalLabScientist.GetMedicalLabScientist1)
	healthProfessional.Get("/get-medicalLabScientist/:id", medicalLabScientist.GetMedicalLabScientistByID)
	healthProfessional.Get("/get-medicalLabScientist-available-slots", medicalLabScientist.GetMedicalLabScientistAvailableTimes)
	healthProfessional.Get("/get-medicalLabScientist-services", medicalLabScientist.GetMedicalLabScientistServices)
	healthProfessional.Post("/add-nurse-appointment", nurse.AddNurseAppointment)
	healthProfessional.Get("/get-nurses", nurse.GetNurses)
	healthProfessional.Get("/get-nurse/:id", nurse.GetNurseByID)
	healthProfessional.Get("/get-nurse-services", nurse.GetNurseServices)
	healthProfessional.Get("/get-nurse-available-slots", nurse.GetNurseAvailableTimes)
	healthProfessional.Post("/add-doctor-appointment", doctorProfession.AddDoctorAppointment)
	healthProfessional.Get("/get-doctor-available-slots", doctorProfession.GetDoctorAvailableTimes)
	healthProfessional.Get("/get-doctors", doctorProfession.GetDoctors)
	healthProfessional.Get("/get-doctor/:id", doctorProfession.GetDoctorProfessionByID)

	appointments := healthFacility.Group("/appointment")
	appointments.Use(jwt, middlewares.CustomerData)
	appointments.Get("/hospital-appointments", hospitals.FetchHospitalAppointmentsWithPagination)
	appointments.Get("/hospital-approved-appointments", hospitals.FetchApprovedHospitalAppointmentsWithPagination)
	// appointments.Get("/hospital-appointment/:id", hospitals.GetHospitalAppointmentByID)
	appointments.Get("/pharmacy-drugs", pharmacy.FetchPharmacyDrugsWithPagination)
	appointments.Get("/pharmacy-approved-drugs", pharmacy.FetchApprovedPharmacyDrugsWithPagination)
	appointments.Get("/fitnessCenter-appointments", fitnessCenter.FetchFitnessCenterAppointmentsWithPagination)
	appointments.Get("/fitnessCenter-approved-appointments", fitnessCenter.FetchApprovedFitnessCenterAppointmentsWithPagination)
	appointments.Get("/fitnessCenter-appointment/:id", fitnessCenter.GetFitnessCenterAppointmentByID)
	appointments.Get("/laboratory-appointments", laboratory.FetchLaboratoryAppointmentsWithPagination)
	appointments.Get("/laboratory-approved-appointments", laboratory.FetchApprovedLaboratoryAppointmentsWithPagination)

	appointment := healthProfessional.Group("/appointment")
	appointments.Use(jwt, middlewares.CustomerData)
	appointment.Get("/doctor-appointments", doctorProfession.FetchDoctorAppointmentsWithPagination)
	appointment.Get("/doctor-approved-appointments", doctorProfession.FetchApprovedDoctorAppointmentsWithPagination)
	appointment.Get("/nurse-appointment/:id", doctorProfession.GetdoctorProfessionAppointmentByID)
	appointment.Get("/nurse-appointments", nurse.FetchNurseAppointmentsWithPagination)
	appointment.Get("/nurse-approved-appointments", nurse.FetchApprovedNurseAppointmentsWithPagination)
	appointment.Get("/nurse-appointment/:id", nurse.GetNurseAppointmentByID)
	appointment.Get("/physiotherapist-appointments", physiotherapist.FetchPhysiotherapistAppointmentsWithPagination)
	appointment.Get("/physiotherapist-approved-appointments", physiotherapist.FetchApprovedPhysiotherapistAppointmentsWithPagination)
	appointment.Get("/physiotherapist-appointment/:id", physiotherapist.GetPhysiotherpistAppointmentByID)
	appointment.Get("/medicalLabScientist-appointments", medicalLabScientist.FetchMedicalLabScientistAppointmentsWithPagination)
	appointment.Get("/medicalLabScientist-approved-appointments", medicalLabScientist.FetchApprovedMedicalLabScientistAppointmentsWithPagination)
	appointment.Get("/medicalLabScientist-appointment/:id", medicalLabScientist.GetMedicalLabScientistAppointmentByID)

	emergencyApis := customer.Group("/emergency")
	emergencyApis.Use(jwt, middlewares.CustomerData)
	emergencyApis.Get("/get-emergency-doctors", emergency.GetEmergencyDoctors)
	emergencyApis.Get("/get-emergency-hospitals", emergency.GetEmergencyHospitals)
	emergencyApis.Post("/add-emergency-doctor", emergency.AddEmergencyDoctor)
	emergencyApis.Post("/add-emergency-hospital", emergency.AddEmergencyHospital)
}
