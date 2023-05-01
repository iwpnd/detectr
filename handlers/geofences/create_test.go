package geofences

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/iwpnd/detectr/database"
	"github.com/iwpnd/detectr/logger"
	"github.com/stretchr/testify/assert"
)

func setupApp() *fiber.App {
	app := fiber.New()

	_ = logger.SetLogLevel("warn")
	lg, _ := logger.New()
	db := database.New()
	RegisterRoutes(app, db, lg)

	return app
}

func TestCreate(t *testing.T) {
	app := setupApp()

	data := []byte(`{"type":"Feature","properties":{},"geometry":{"type":"Polygon","coordinates":[[[13.3967096231641,52.47425410999395],[13.3967096231641,52.4680479999262],[13.413318577304466,52.4680479999262],[13.413318577304466,52.47425410999395],[13.3967096231641,52.47425410999395]]]}}`)

	type tcase struct {
		Body         []byte
		ContentType  string
		ExpectedCode int
	}

	tests := map[string]tcase{
		"test application/json": {
			Body:         data,
			ContentType:  "application/json",
			ExpectedCode: 201,
		},
		"test application/geo+json": {
			Body:         data,
			ContentType:  "application/geo+json",
			ExpectedCode: 422,
		},
		"test empty geofence": {
			Body:         []byte(``),
			ContentType:  "application/json",
			ExpectedCode: 422,
		},
		"test faulty geofence": {
			Body:         []byte(`{"foo":"bar"}`),
			ContentType:  "application/json",
			ExpectedCode: 422,
		},
	}

	for _, test := range tests {

		r := httptest.NewRequest(
			"POST",
			"/geofence",
			bytes.NewBuffer(test.Body),
		)
		r.Header.Add("Content-Type", test.ContentType)

		resp, _ := app.Test(r, -1)
		assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	}
}
