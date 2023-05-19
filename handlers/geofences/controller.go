package geofences

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iwpnd/detectr/database"
	"go.uber.org/zap"
)

type handler struct {
	DB     database.Datastore
	Logger *zap.Logger
}

// RegisterRoutes to register geofence routes with the fiber app
func RegisterRoutes(app *fiber.App, db database.Datastore, logger *zap.Logger) {
	h := &handler{
		DB:     db,
		Logger: logger,
	}

	routes := app.Group("geofence")
	routes.Post("/", h.CreateFence)
}
