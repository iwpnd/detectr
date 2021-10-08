package main

import (
	"github.com/buckhx/diglet/geo"
	"github.com/gofiber/fiber/v2"
	"github.com/iwpnd/detectr/detectr"
	// "github.com/iwpnd/fiber-key-auth"
)

func main() {
	col := detectr.NewCollection()

	o := `{type":"Feature","properties":{"spec":"reduced-speed"},"geometry":{"type":"Polygon","coordinates":[[[12.8814697265625,52.26815737376817],[13.809814453125,52.26815737376817],[13.809814453125,52.76289173758374],[12.8814697265625,52.76289173758374],[12.8814697265625,52.26815737376817]]]}}`

	g, err := geo.UnmarshalGeojsonFeature(o)

	if err != nil {
		return
	}

	feature, err := geo.GeojsonFeatureAdapter(g)

	if err != nil {
		return
	}

	col.Insert(feature)

	app := fiber.New()
	// app.Use(keyauth.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	app.Get("/index", func(c *fiber.Ctx) error {
		p := geo.Coordinate{Lon: 13.5296630859375, Lat: 52.43926935464697}

		all := col.Contains(p)

		return c.JSON(all)

	})

	app.Listen(":3000")
}
