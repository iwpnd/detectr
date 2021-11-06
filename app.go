package main

import (
	"github.com/iwpnd/detectr/database"
	"github.com/iwpnd/detectr/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/iwpnd/fiber-key-auth"
)

func main() {
	f := database.Get()
	f.LoadFromPath("bin/test.geojson")

	app := fiber.New()
	app.Use(keyauth.New())

	app.Get("/healthz", handlers.GetHealthz)
	app.Post("/fence", handlers.PostFence)
	app.Post("/location", handlers.PostLocation)

	app.Listen(":3000")
}
