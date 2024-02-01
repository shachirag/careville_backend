package middlewares

import (
	adminMiddleware "careville_backend/dto/admin/middleware"

	"github.com/gofiber/fiber/v2"
)

func AdminData(c *fiber.Ctx) error {
	_, err := adminMiddleware.SetAdminMiddlewareData(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseDto{
			Status:  false,
			Message: "Failed to set admin middleware data: " + err.Error(),
		})
	}

	return c.Next()
}
