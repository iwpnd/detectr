package main

import (
	"github.com/iwpnd/detectr/collection"
	"github.com/iwpnd/detectr/models"
	"github.com/iwpnd/detectr/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/tidwall/geojson"
	"github.com/tidwall/geojson/geometry"
	// "github.com/iwpnd/fiber-key-auth"
)

type resp struct {
	objects int
}

func main() {
	col := collection.NewCollection()

	o := `{"type":"Feature","properties":{"spec":"reduced-speed"},"geometry":{"type":"Polygon","coordinates":[[[12.8814697265625,52.26815737376817],[13.809814453125,52.26815737376817],[13.809814453125,52.76289173758374],[12.8814697265625,52.76289173758374],[12.8814697265625,52.26815737376817]]]}}`

	g, err := geojson.Parse(o, nil)

	if err != nil {
		return
	}

	col.Insert(g)

	app := fiber.New()
	// app.Use(keyauth.New())

	app.Get("/healtz", func(c *fiber.Ctx) error {
		return c.SendString("")
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

		items := col.Intersects(p)
		return c.JSON(items)
	})

	app.Listen(":3000")
}
