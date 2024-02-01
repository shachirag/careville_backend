package router

import (
	"careville_backend/handlers"
	customerAuth "careville_backend/handlers/customer/customerAuthentication"
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
}
