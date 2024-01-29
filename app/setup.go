package app

import (
	"careville_backend/config"
	"careville_backend/database"
	"careville_backend/router"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupAndRunApp() error {
	// loaded the env config
	err := config.LoadENV()
	if err != nil {
		return err
	}

	// created the database
	err = database.StartMongoDB()
	if err != nil {
		return err
	}

	// defer closing the database
	defer database.CloseMongoDB()

	// Setup s3
	err = database.SetupAWSClient()
	if err != nil {
		return err
	}

	// new fiber instance
	app := fiber.New(fiber.Config{
		BodyLimit: 16 * 1024 * 1204,
	})

	// middlewares
	app.Use(recover.New())
	app.Use(cors.New(cors.ConfigDefault))
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} ${latency}\n",
	}))

	// set up the Provider routes
	router.ProviderSetupsRoutes(app)

	// set up the admin routes
	router.AdminSetupsRoutes(app)

	// Setup Swagger
	config.AddSwaggerRoutes(app)

	// create the address for listening
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":5500"
	}

	// listen on the address
	err = app.Listen(addr)
	if err != nil {
		return err
	}

	return nil
}
