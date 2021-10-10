package handlers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/iwpnd/detectr/fences"
	"github.com/iwpnd/detectr/models"
	"github.com/iwpnd/detectr/validation"
	"github.com/tidwall/geojson"
	"github.com/tidwall/geojson/geometry"
	"time"
)

type FenceResponse struct {
	Elapsed string           `json:"elapsed"`
	Request models.Location  `json:"request"`
	Fences  []geojson.Object `json:"fences"`
}

type Response struct {
	Data interface{} `json:"data"`
}

func PostLocation(c *fiber.Ctx) error {
	start := time.Now()
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
	elapsed := time.Since(start)

	fr := &Response{
		Data: FenceResponse{
			Elapsed: fmt.Sprint(elapsed),
			Request: *location,
			Fences:  matches,
		},
	}

	return c.JSON(&fr)
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

	resp := &Response{Data: d}
	return c.JSON(&resp)
}

func GetHealthz(c *fiber.Ctx) error {
	return c.SendStatus(200)
}
