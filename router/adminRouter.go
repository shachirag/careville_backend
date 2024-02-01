package router

import (
	"careville_backend/handlers"
	adminAuth "careville_backend/handlers/admin/adminAuthorization"
	dashboard "careville_backend/handlers/admin/dashboard"
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

	adminProfile := admin.Group("/profile")
	adminProfile.Use(jwt, middlewares.AdminData)
	adminProfile.Get("/get-admin-info", adminAuth.FetchAdminById)
	adminProfile.Put("/change-password", adminAuth.ChangeAdminPassword)
	adminProfile.Put("/update-admin-info", adminAuth.UpdateAdmin)
	adminProfile.Put("/change-status/:id", adminAuth.ChangeStatus)

	// dashboard
	adminProfile.Get("/total-counts", jwt, dashboard.CountAll)

}
