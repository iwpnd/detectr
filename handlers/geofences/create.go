package geofences

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iwpnd/detectr/models"
	"github.com/tidwall/geojson"
)

func (h *handler) CreateFence(c *fiber.Ctx) error {
	d, err := geojson.Parse(string(c.Body()), nil)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	h.DB.Create(d)

	r := &models.Response{Data: d}
	return c.JSON(&r)
}
