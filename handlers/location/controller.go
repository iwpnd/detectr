package location

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

	routes := app.Group("location")
	routes.Post("/", h.PostLocation)
}
