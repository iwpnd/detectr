package location

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/iwpnd/detectr/models"
	"github.com/iwpnd/detectr/validation"
	"github.com/tidwall/geojson"
	"github.com/tidwall/geojson/geometry"
)

func (h *handler) PostLocation(c *fiber.Ctx) error {
	start := time.Now()
	l := new(models.Location)

	if err := c.BodyParser(l); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

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
	elapsed := time.Since(start)

	r := &models.Response{
		Data: models.GeofenceResponse{
			Elapsed: fmt.Sprint(elapsed),
			Request: *l,
			Matches: matches,
		},
	}

	return c.JSON(&r)
}
