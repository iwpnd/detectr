package main

import (
	"github.com/iwpnd/detectr/collection"
	"github.com/iwpnd/detectr/models"
	"github.com/iwpnd/detectr/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/iwpnd/fiber-key-auth"
	"github.com/tidwall/geojson"
	"github.com/tidwall/geojson/geometry"
)

func main() {
	col := collection.NewCollection()
	col.LoadFromPath("bin/test.geojson")

	app := fiber.New()
	app.Use(keyauth.New())

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Post("/fence", func(c *fiber.Ctx) error {
		d, err := geojson.Parse(string(c.Body()), nil)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		col.Create(d)
		return c.JSON(d)
	})

	app.Post("/location", func(c *fiber.Ctx) error {
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

		matches := col.Intersects(p)
		return c.JSON(geojson.NewFeatureCollection(matches))
	})

	app.Listen(":3000")
}
