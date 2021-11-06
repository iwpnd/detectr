package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/iwpnd/detectr/database"
	"github.com/iwpnd/detectr/models"
	"github.com/iwpnd/detectr/validation"
	"github.com/tidwall/geojson"
	"github.com/tidwall/geojson/geometry"
	"time"
)

type FenceResponse struct {
	Elapsed string           `json:"elapsed"`
	Request models.Location  `json:"request"`
	Matches []geojson.Object `json:"matches"`
}

type Response struct {
	Data interface{} `json:"data"`
}

func PostLocation(c *fiber.Ctx) error {
	start := time.Now()
	f := database.Get()
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

	matches := f.Intersects(p)
	elapsed := time.Since(start)

	r := &Response{
		Data: FenceResponse{
			Elapsed: fmt.Sprint(elapsed),
			Request: *l,
			Matches: matches,
		},
	}

	return c.JSON(&r)
}

func PostFence(c *fiber.Ctx) error {
	f := database.Get()

	d, err := geojson.Parse(string(c.Body()), nil)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	f.Create(d)

	r := &Response{Data: d}
	return c.JSON(&r)
}

func GetHealthz(c *fiber.Ctx) error {
	return c.SendStatus(200)
}
