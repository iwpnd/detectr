package geofences

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/iwpnd/detectr/errors"
	"github.com/iwpnd/detectr/models"
	geojson "github.com/paulmach/go.geojson"
	"go.uber.org/zap"
)

// CreateFence ...
func (h *handler) CreateFence(c *fiber.Ctx) error {
	start := time.Now()

	if string(c.Request().Header.ContentType()) != "application/json" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors.NewRequestError(&errors.ErrRequestError{
			Status: fiber.StatusUnprocessableEntity,
			Detail: "Content-Type must be 'application/json'",
		}))
	}

	d, err := geojson.UnmarshalFeature(c.Body())
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors.NewRequestError(&errors.ErrRequestError{Status: fiber.StatusUnprocessableEntity,
			Detail: err.Error(),
		}))
	}

	err = h.DB.Create(d)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors.NewRequestError(&errors.ErrRequestError{
			Status: fiber.StatusUnprocessableEntity,
			Detail: err.Error(),
		}))
	}

	r := &models.Response{Data: d}

	elapsed := time.Since(start)
	h.Logger.Debug("Created geofence", zap.Any("data", &d), zap.String("elapsed", fmt.Sprint(elapsed)))

	return c.Status(201).JSON(&r)
}
