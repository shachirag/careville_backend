package middlewares

import (
	providerMiddleware "careville_backend/dto/provider/middleware"
	"github.com/gofiber/fiber/v2"
)

type ResponseDto struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func ProviderData(c *fiber.Ctx) error {
	_, err := providerMiddleware.SetProviderMiddlewareData(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseDto{
			Status:  false,
			Message: "Failed to set provider middleware data: " + err.Error(),
		})
	}

	return c.Next()
}
