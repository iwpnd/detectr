package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iwpnd/detectr/fences"
	"github.com/iwpnd/detectr/models"
	"github.com/iwpnd/detectr/validation"
	"github.com/tidwall/geojson"
	"github.com/tidwall/geojson/geometry"
)

func PostLocation(c *fiber.Ctx) error {
	f := fences.Get()
	location := new(models.Location)

	if err := c.BodyParser(location); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := validation.ValidateStruct(*location)
	if errors != nil {
		return c.JSON(errors)
	}

	p := geojson.NewPoint(
		geometry.Point{
			X: location.Lng,
			Y: location.Lat,
		},
	)

	matches := f.Intersects(p)
	return c.JSON(geojson.NewFeatureCollection(matches))
}

func PostFence(c *fiber.Ctx) error {
	f := fences.Get()

	d, err := geojson.Parse(string(c.Body()), nil)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	f.Create(d)
	return c.JSON(d)
}

func GetHealthz(c *fiber.Ctx) error {
	return c.SendStatus(200)
}
