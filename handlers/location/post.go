package location

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/iwpnd/detectr/models"
	"github.com/iwpnd/detectr/validation"
	"github.com/tidwall/geojson"
	"github.com/tidwall/geojson/geometry"
	"go.uber.org/zap"
)

func (h *handler) PostLocation(c *fiber.Ctx) error {
	start := time.Now()
	l := new(models.Location)

	if err := c.BodyParser(l); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	h.Logger.Debug("Received Location",
		zap.Float64("latitude", l.Lat),
		zap.Float64("longitude", l.Lng),
	)

	errors := validation.ValidateStruct(*l)
	if errors != nil {
		return c.JSON(errors)
	}

	p := geojson.NewPoint(
		geometry.Point{
			X: l.Lng,
			Y: l.Lat,
		},
	)

	matches := h.DB.Intersects(p)
	elapsed := fmt.Sprint(time.Since(start))

	r := &models.Response{
		Data: models.GeofenceResponse{
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
