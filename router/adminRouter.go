package router

import (
	"careville_backend/handlers"
	adminAuth "careville_backend/handlers/admin/adminAuthorization"
	dashboard "careville_backend/handlers/admin/dashboard"

	"github.com/gofiber/fiber/v2"
)

func AdminSetupsRoutes(app *fiber.App) {

	app.Static("/", ".puplic")

	/* ---------- HEALTH ---------- */
	app.Get("/health", handlers.HealthCheck)

	/* ---------- Protected Routes ----- */
	// secret := os.Getenv("JWT_SECRET_KEY")
	// jwt := middlewares.NewAuthMiddleware(secret)

	// provider authentication
	admin := app.Group("/admin")
	admin.Post("/login", adminAuth.LoginAdmin)
	admin.Post("/verify-otp-for-login", adminAuth.VerifyOtpForLogin)
	admin.Get("/get-admin-info/:adminId", adminAuth.FetchAdminById)
	admin.Put("/change-password/:adminId", adminAuth.ChangeAdminPassword)
	admin.Put("/update-admin-info/:adminId", adminAuth.UpdateAdmin)
	admin.Put("/change-status/:id", adminAuth.ChangeStatus)

	// dashboard
	admin.Get("/total-counts", dashboard.CountAll)

}
