package location

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iwpnd/detectr/database"
	"go.uber.org/zap"
)

type handler struct {
	DB     *database.Database
	Logger *zap.Logger
}

func RegisterRoutes(app *fiber.App, db *database.Database, logger *zap.Logger) {
	h := &handler{
		DB:     db,
		Logger: logger,
	}

	routes := app.Group("location")
	routes.Post("/", h.PostLocation)
}
