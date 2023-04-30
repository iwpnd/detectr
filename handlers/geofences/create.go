package geofences

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/iwpnd/detectr/models"
	"github.com/tidwall/geojson"
	"go.uber.org/zap"
)

func (h *handler) CreateFence(c *fiber.Ctx) error {
	start := time.Now()
	d, err := geojson.Parse(string(c.Body()), nil)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	h.DB.Create(d)

	r := &models.Response{Data: d}

	elapsed := time.Since(start)
	h.Logger.Debug("Created geofence", zap.Any("data", &d), zap.String("elapsed", fmt.Sprint(elapsed)))

	return c.JSON(&r)
}
