package location

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/iwpnd/detectr/errors"
	"github.com/iwpnd/detectr/models"
	"github.com/iwpnd/detectr/validation"
	"github.com/tidwall/geojson"
	"github.com/tidwall/geojson/geometry"
	"go.uber.org/zap"
)

func (h *handler) PostLocation(c *fiber.Ctx) error {
	start := time.Now()
	l := new(models.Location)

	if string(c.Request().Header.ContentType()) != "application/json" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors.NewRequestError(&errors.ErrRequestError{
			Status: fiber.StatusUnprocessableEntity,
			Detail: "Content-type must be 'application/json'",
		}))
	}

	if err := c.BodyParser(l); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			errors.NewRequestError(&errors.ErrRequestError{Detail: err.Error()}),
		)
	}

	h.Logger.Debug("Received Location",
		zap.Float64("latitude", l.Lat),
		zap.Float64("longitude", l.Lng),
	)

	errs := validation.ValidateStruct(*l)
	if errs != nil {
		return c.Status(400).JSON(errors.NewRequestError(errs...))
	}

	p := geojson.NewPoint(
		geometry.Point{
			X: l.Lng,
			Y: l.Lat,
		},
	)

	matches := h.DB.Intersects(p)
	elapsed := fmt.Sprint(time.Since(start))

	r := &models.LocationResponse{
		Data: models.LocationResponsePayload{
			Elapsed: elapsed,
			Request: *l,
			Matches: matches,
		},
	}

	h.Logger.Debug("Processed Location",
		zap.Float64("latitude", l.Lat),
		zap.Float64("longitude", l.Lng),
		zap.String("elapsed", elapsed),
		zap.Any("matches", &matches),
	)

	return c.JSON(&r)
}
