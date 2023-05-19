package memory

import (
	"testing"

	geojson "github.com/paulmach/go.geojson"
	"github.com/stretchr/testify/assert"
)

func setupDatabase() *Memory {
	db := New()

	return db
}

func TestCreate(t *testing.T) {
	db := setupDatabase()
	defer db.Truncate()

	expected := []byte(`{"id":"foobar","type":"Feature","geometry":{"type":"Polygon","coordinates":[[[13.3967096231641,52.47425410999395],[13.3967096231641,52.4680479999262],[13.413318577304466,52.4680479999262],[13.413318577304466,52.47425410999395],[13.3967096231641,52.47425410999395]]]},"properties":{"id":"foobar"}}`)

	f, err := geojson.UnmarshalFeature(expected)
	if err != nil {
		t.Fatal("failed to unmarshal feature: ", err)
	}

	err = db.Create(f)
	if err != nil {
		t.Fatal("failed to create feature")
	}

	p := []float64{13.40532627661105, 52.471361312503575}

	matches := db.Intersects(p)
	got, _ := matches[0].MarshalJSON()

	assert.Equal(t, string(expected), string(got))
}

func TestCreateGenerateID(t *testing.T) {
	db := setupDatabase()
	defer db.Truncate()

	data := []byte(`{"type":"Feature","properties":{},"geometry":{"type":"Polygon","coordinates":[[[13.3967096231641,52.47425410999395],[13.3967096231641,52.4680479999262],[13.413318577304466,52.4680479999262],[13.413318577304466,52.47425410999395],[13.3967096231641,52.47425410999395]]]}}`)

	f, err := geojson.UnmarshalFeature(data)
	if err != nil {
		t.Fatal("failed to unmarshal feature: ", err)
	}

	assert.Nil(t, f.ID)

	err = db.Create(f)
	if err != nil {
		t.Fatal("failed to create feature")
	}

	p := []float64{13.40532627661105, 52.471361312503575}

	matches := db.Intersects(p)

	assert.NotNil(t, matches[0].ID)
}

func TestCreateFailed(t *testing.T) {
	db := setupDatabase()
	defer db.Truncate()

	type tcase struct {
		Data          []byte
		expectedError string
	}

	tests := map[string]tcase{
		"test invalid geometry type": {
			Data:          []byte(`{"type":"Feature","properties":{},"geometry":{"type":"Point","coordinates":[1,1]}}`),
			expectedError: "unsupported geometry type: Point",
		},
		"test faulty geofence": {
			Data:          []byte(`{"foo":"bar"}`),
			expectedError: "empty geometry",
		},
	}

	for _, test := range tests {
		f, err := geojson.UnmarshalFeature(test.Data)
		if err != nil {
			t.Fatal("failed to unmarshal feature: ", err)
		}
		err = db.Create(f)
		if test.expectedError != "" && err != nil {
			assert.Equal(t, test.expectedError, err.Error())
		}
	}
}

func TestTruncate(t *testing.T) {
	db := setupDatabase()

	data := []byte(`{"type":"Feature","properties":{},"geometry":{"type":"Polygon","coordinates":[[[13.3967096231641,52.47425410999395],[13.3967096231641,52.4680479999262],[13.413318577304466,52.4680479999262],[13.413318577304466,52.47425410999395],[13.3967096231641,52.47425410999395]]]}}`)

	f, err := geojson.UnmarshalFeature(data)
	if err != nil {
		t.Fatal("failed to unmarshal feature: ", err)
	}
	err = db.Create(f)
	if err != nil {
		t.Fatal("failed to create feature")
	}

	p := []float64{13.40532627661105, 52.471361312503575}

	assert.Equal(t, 1, len(db.Intersects(p)))

	db.Truncate()
	assert.Equal(t, 0, db.Count())
	assert.Equal(t, 0, len(db.Intersects(p)))
}

func TestDelete(t *testing.T) {
	db := setupDatabase()
	defer db.Truncate()

	data := []byte(`{"type":"Feature","properties":{},"geometry":{"type":"Polygon","coordinates":[[[13.3967096231641,52.47425410999395],[13.3967096231641,52.4680479999262],[13.413318577304466,52.4680479999262],[13.413318577304466,52.47425410999395],[13.3967096231641,52.47425410999395]]]}}`)

	f, err := geojson.UnmarshalFeature(data)
	if err != nil {
		t.Fatal("failed to unmarshal feature: ", err)
	}
	err = db.Create(f)
	if err != nil {
		t.Fatal("failed to create feature")
	}

	p := []float64{13.40532627661105, 52.471361312503575}

	assert.Equal(t, 1, db.Count())
	assert.Equal(t, 1, len(db.Intersects(p)))

	db.Delete(f)

	assert.Equal(t, 0, db.Count())
	assert.Equal(t, 0, len(db.Intersects(p)))
}
