# detectr

a minimal geofencing application build with [gofiber](https://github.com/gofiber/fiber),
[rtree](https://github.com/tidwall/rtree) and [geoindex](https://github.com/tidwall/geoindex).

Bring your own data, pass a location and receive the geofences that your location
is currently in.

## Installation

### cli

```bash
go install github.com/iwpnd/detectr/cmd/detectr@latest
```

```bash
➜ detectr --help
NAME:
   detectr - geofence application

USAGE:
   detectr [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --port value       define port (default: 3000)
   --require-key      use keyauth (default: false)
   --log-level value  set loglevel (default: "error")
   --data value       path to dataset to load with app
   --help, -h         show help (default: false)
```

## Usage

```bash
# start up

detectr --data my_featurecollection.geojson
```

```bash
# send a location

curl -X POST -d "lat=13.4042367&lng=52.473091" http://localhost:3000/location

>> {
  "data": {
    "elapsed": "150.75µs",
    "request": {
      "lat": 52.25,
      "lng": 13.37
    },
    "matches": [
        {
          "type": "Feature",
          "properties": {"id":"my geofence"},
          "geometry": {
            "coordinates": [
              [
                [
                  13.41493887975497,
                  52.47961028115867
                ],
                [
                  13.393534522441712,
                  52.47961028115867
                ],
                [
                  13.393534522441712,
                  52.466572160399664
                ],
                [
                  13.41493887975497,
                  52.466572160399664
                ],
                [
                  13.41493887975497,
                  52.47961028115867
                ]
              ]
            ],
            "type": "Polygon"
          }
        }
    ]
  }
}
```

## License

MIT

## Acknowledgement

- Triggered by this article on the UBER engineering blog
  about [UBERs highest query per second service](https://www.uber.com/en-DE/blog/go-geofence-highest-query-per-second-service/).
- Heavily inspired by [Tile38](https://github.com/tidwall/tile38) and build because
  I wanted to understand how Tile38 works under the hood.

## Maintainer

Benjamin Ramser - [@iwpnd](https://github.com/iwpnd)

Project Link: [https://github.com/iwpnd/detectr](https://github.com/iwpnd/detectr)
