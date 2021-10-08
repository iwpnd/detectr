package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/iwpnd/detectr/detectr"
	"github.com/tidwall/geojson"
	"github.com/tidwall/geojson/geometry"
	// "github.com/iwpnd/fiber-key-auth"
)

type resp struct {
	objects int
}

func main() {
	col := detectr.NewCollection()

	o := `{type":"Feature","properties":{"spec":"reduced-speed"},"geometry":{"type":"Polygon","coordinates":[[[12.8814697265625,52.26815737376817],[13.809814453125,52.26815737376817],[13.809814453125,52.76289173758374],[12.8814697265625,52.76289173758374],[12.8814697265625,52.26815737376817]]]}}`

	g, err := geojson.Parse(o, nil)

	if err != nil {
		return
	}

	col.Insert(g)

	app := fiber.New()
	// app.Use(keyauth.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Get("/index", func(c *fiber.Ctx) error {

		p := geojson.NewPoint(geometry.Point{X: 13.809814453125, Y: 52.26815737376817})
		fmt.Print(p)

		return c.JSON(resp{objects: col.Count()})

	})

	app.Listen(":3000")
}
