package main

import (
	"log"
	"os"

	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/iwpnd/detectr/database"
	"github.com/iwpnd/detectr/handlers/geofences"
	"github.com/iwpnd/detectr/handlers/location"
	"github.com/iwpnd/detectr/logger"

	"github.com/gofiber/fiber/v2"
	keyauth "github.com/iwpnd/fiber-key-auth"
)

func startDetectr(ctx *cli.Context) error {
	port := ctx.Int64("port")
	datapath := ctx.String("data")
	loglevel := ctx.String("log-level")
	requirekey := ctx.Bool("require-key")

	db := database.New()

	c := fiber.Config{
		AppName: "detectr",
	}
	app := fiber.New(c)
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	if datapath != "" {
		err := db.LoadFromPath(datapath)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	if loglevel != "" {
		err := logger.SetLogLevel(loglevel)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	if requirekey {
		app.Use(keyauth.New(keyauth.WithStructuredErrorMsg()))
	}

	logr, err := logger.New()
	if err != nil {
		fmt.Println(err)
	}

	location.RegisterRoutes(app, db, logr)
	geofences.RegisterRoutes(app, db, logr)

	err = app.Listen(fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

var withPort cli.Int64Flag
var withKeyAuth cli.BoolFlag
var withLogLevel cli.StringFlag
var withDataset cli.StringFlag

func init() {
	withPort = cli.Int64Flag{
		Name:     "port",
		Usage:    "define port",
		Value:    3000,
		Required: true,
	}
	withKeyAuth = cli.BoolFlag{
		Name:     "require-key",
		Usage:    "use keyauth",
		Value:    false,
		Required: false,
	}
	withLogLevel = cli.StringFlag{
		Name:     "log-level",
		Usage:    "set loglevel",
		Value:    "error",
		Required: false,
	}
	withDataset = cli.StringFlag{
		Name:     "data",
		Usage:    "path to dataset to load with app",
		Required: false,
	}
}

func main() {
	app := &cli.App{
		Name:   "detectr",
		Usage:  "geofence application",
		Action: startDetectr,
		Flags: []cli.Flag{
			&withPort,
			&withKeyAuth,
			&withLogLevel,
			&withDataset,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
