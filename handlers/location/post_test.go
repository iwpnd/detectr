package location

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/iwpnd/detectr/database"
	"github.com/iwpnd/detectr/logger"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/geojson"
)

func locationRequest(lat float64, lng float64) []byte {
	return []byte(fmt.Sprintf(`{"lat":%v,"lng":%v}`, lat, lng))
}

func setupApp() (*fiber.App, geojson.Object) {
	app := fiber.New()

	data, _ := geojson.Parse(string(
		[]byte(`{"type":"Feature","properties":{},"geometry":{"type":"Polygon","coordinates":[[[-3.0988311767578125,40.837710162420045],[-3.121490478515625,40.820045086716505],[-3.0978012084960938,40.80237530523985],[-3.0754852294921875,40.8210843390845],[-3.0988311767578125,40.837710162420045]],[[-3.0988311767578125,40.82783908257347],[-3.1098175048828125,40.820045086716505],[-3.0988311767578125,40.81147063339219],[-3.086471557617187,40.820304901335035],[-3.0988311767578125,40.82783908257347]]]}}`),
	), nil)

	lg, _ := logger.New()

	db := database.New()
	db.Create(data)

	RegisterRoutes(app, db, lg)

	return app, data
}

func TestLocation(t *testing.T) {
	app, _ := setupApp()

	type tcase struct {
		Body            []byte
		ContentType     string
		ExpectedCode    int
		ExpectedMatches string
	}

	fmt.Print(string(locationRequest(1, 1)))

	tests := map[string]tcase{
		"south-east outside polygon, in bbox": {
			Body:         locationRequest(40.80809251416925, -3.0816650390625),
			ContentType:  "application/json",
			ExpectedCode: 200,
		},
		"south-east inside polygon, inside bbox": {
			Body:         locationRequest(40.81497849824719, -3.0878448486328125),
			ContentType:  "application/json",
			ExpectedCode: 200,
		},
		"south-east outside polygon, outside bbox": {
			Body:         locationRequest(40.800945926051526, -3.0713653564453125),
			ContentType:  "application/json",
			ExpectedCode: 200,
		},
		"south outside polygon, outside bbox": {
			Body:         locationRequest(40.79769722250925, -3.0978012084960938),
			ContentType:  "application/json",
			ExpectedCode: 200,
		},
		"south inside polygon, inside bbox": {
			Body:         locationRequest(40.8067931917519, -3.098316192626953),
			ContentType:  "application/json",
			ExpectedCode: 200,
		},
		"south-west outside polygon, inside bbox": {
			Body:         locationRequest(40.807702720115294, -3.116168975830078),
			ContentType:  "application/json",
			ExpectedCode: 200,
		},
		"south-west outside polygon, outside bbox": {
			Body:         locationRequest(40.80068603561921, -3.1250953674316406),
			ContentType:  "application/json",
			ExpectedCode: 200,
		},
		"south-west inside polygon, inside bbox": {
			Body:         locationRequest(40.814198988751876, -3.10810089111328),
			ContentType:  "application/json",
			ExpectedCode: 200,
		},
		"west outside polygon, outside bbox": {
			Body:         locationRequest(40.8197852710803, -3.1266403198242188),
			ContentType:  "application/json",
			ExpectedCode: 200,
		},
		"west inside polygon, inside bbox": {
			Body:         locationRequest(40.82017499415298, -3.1141090393066406),
			ContentType:  "application/json",
			ExpectedCode: 200,
		},
		"north-west inside polygon, inside bbox": {
			Body:         locationRequest(40.82667004158603, -3.1070709228515625),
			ContentType:  "application/json",
			ExpectedCode: 200,
		},
		"north-west outside polygon, inside bbox": {
			Body:         locationRequest(40.83199550584334, -3.1141090393066406),
			ContentType:  "application/json",
			ExpectedCode: 200,
		},
		"north inside polygon, inside bbox": {
			Body:         locationRequest(40.83264492344398, -3.0988311767578125),
			ContentType:  "application/json",
			ExpectedCode: 200,
		},
		"north outside polygon, outside bbox": {
			Body:         locationRequest(40.8425152878029, -3.0988311767578125),
			ContentType:  "application/json",
			ExpectedCode: 200,
		},
		"north-east inside polygon, inside bbox": {
			Body:         locationRequest(40.826799936046804, -3.0895614624023438),
			ContentType:  "application/json",
			ExpectedCode: 200,
		},
		"north-east outside polygon, inside bbox": {
			Body:         locationRequest(40.83160585222969, -3.0816650390625),
			ContentType:  "application/json",
			ExpectedCode: 200,
		},
		"north-east outside polygon, outside bbox": {
			Body:         locationRequest(40.84147637129013, -3.07016372680664),
			ContentType:  "application/json",
			ExpectedCode: 200,
		},
		"east inside polygon, inside bbox": {
			Body:         locationRequest(40.82056471493589, -3.080635070800781),
			ContentType:  "application/json",
			ExpectedCode: 200,
		},
		"east outside polygon, outside bbox": {
			Body:         locationRequest(40.8210843390845, -3.069477081298828),
			ContentType:  "application/json",
			ExpectedCode: 200,
		},
		"east outside polygon but in hole, inside bbox": {
			Body:         locationRequest(40.81874599835864, -3.098487854003906),
			ContentType:  "application/json",
			ExpectedCode: 200,
		},
		"fails because wrong content type": {
			Body:         locationRequest(40.81874599835864, -3.098487854003906),
			ContentType:  "application/geo+json",
			ExpectedCode: 422,
		},
		"fails because bad input": {
			Body:         locationRequest(1140.81874599835864, -3.098487854003906),
			ContentType:  "application/json",
			ExpectedCode: 400,
		},
	}

	for _, test := range tests {
		req := httptest.NewRequest(
			"POST",
			"/location",
			bytes.NewBuffer(test.Body),
		)
		req.Header.Add("Content-Type", test.ContentType)

		resp, _ := app.Test(req, -1)

		assert.Equal(t, test.ExpectedCode, resp.StatusCode)
	}
}
