package middlewares

import (
	customerMiddleware "careville_backend/dto/customer/middleware"

	"github.com/gofiber/fiber/v2"
)

func CustomerData(c *fiber.Ctx) error {
	_, err := customerMiddleware.SetCustomerMiddlewareData(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseDto{
			Status:  false,
			Message: "Failed to set customer middleware data: " + err.Error(),
		})
	}

	return c.Next()
}
