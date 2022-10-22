package geofences

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iwpnd/detectr/database"
)

type handler struct {
	DB *database.Database
}

func RegisterRoutes(app *fiber.App, db *database.Database) {
	h := &handler{
		DB: db,
	}

	routes := app.Group("geofence")
	routes.Post("/", h.CreateFence)
}
