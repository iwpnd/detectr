package main

import (
	"fmt"

	"github.com/iwpnd/detectr/database"
	"github.com/iwpnd/detectr/handlers/geofences"
	"github.com/iwpnd/detectr/handlers/location"
	"github.com/iwpnd/detectr/logger"

	"github.com/gofiber/fiber/v2"
	keyauth "github.com/iwpnd/fiber-key-auth"
)

func main() {
	db := database.New()
	db.LoadFromPath("bin/test.geojson")

	app := fiber.New()
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	logger, err := logger.New()
	if err != nil {
		fmt.Println(err)
	}

	app.Use(keyauth.New(keyauth.WithStructuredErrorMsg()))

	location.RegisterRoutes(app, db, logger)
	geofences.RegisterRoutes(app, db, logger)

	app.Listen(":3000")
}
