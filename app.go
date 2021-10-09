package main

import (
	"fmt"
	"github.com/iwpnd/detectr/collection"
	"github.com/iwpnd/detectr/models"
	"github.com/iwpnd/detectr/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/tidwall/geojson"
	"github.com/tidwall/geojson/geometry"
	// "github.com/iwpnd/fiber-key-auth"
)

func main() {
	col := collection.NewCollection()

	err := col.LoadFromPath("bin/test.geojson")

	if err != nil {
		fmt.Print("Could not load from file")
	}

	fmt.Printf("Successfully inserted %v items to the collection", col.Count())

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

		matches := col.Intersects(p)
		return c.JSON(geojson.NewFeatureCollection(matches))
	})

	app.Listen(":3000")
}
