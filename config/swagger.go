package config

import (
	_ "careville_backend/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func AddSwaggerRoutes(app *fiber.App) {
	app.Get("/docs/*", swagger.HandlerDefault)
}
