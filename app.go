package main

import (
	"github.com/iwpnd/detectr/database"
	"github.com/iwpnd/detectr/handlers/geofences"
	"github.com/iwpnd/detectr/handlers/location"

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

	app.Use(keyauth.New())

	location.RegisterRoutes(app, db)
	geofences.RegisterRoutes(app, db)

	app.Listen(":3000")
}
