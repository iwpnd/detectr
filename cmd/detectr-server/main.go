package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/iwpnd/detectr/database"
	"github.com/iwpnd/detectr/errors"
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

	if loglevel != "" {
		err := logger.SetLogLevel(loglevel)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}

	l, err := logger.New()

	if err != nil {
		fmt.Println(err)
	}

	db := database.New()

	c := fiber.Config{
		AppName:               "detectr",
		DisableStartupMessage: true,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(
				fiber.StatusInternalServerError).JSON(
				errors.NewRequestError(
					&errors.ErrRequestError{
						Status: fiber.StatusInternalServerError,
						Detail: err.Error()}))
		},
	}

	app := fiber.New(c)
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	l.Info("Starting detectr...")

	if port == 0 {
		port = 3000
	}

	if datapath != "" {
		start := time.Now()
		l.Info(fmt.Sprintf("Ingesting data from local file path: %v", datapath))
		err := db.LoadFromPath(datapath)
		if err != nil {
			log.Fatal(err)
			return err
		}
		elapsed := fmt.Sprint(time.Since(start))
		l.Info(fmt.Sprintf("Done ingesting data from local file path: %v. Took: %v", datapath, elapsed))
	}

	if requirekey {
		app.Use(keyauth.New(keyauth.WithStructuredErrorMsg()))
	}

	location.RegisterRoutes(app, db, l)
	geofences.RegisterRoutes(app, db, l)

	go func() {
		l.Info(fmt.Sprintf("Detectr listening on port :%v", port))
		if err := app.Listen(fmt.Sprintf(":%v", port)); err != nil {
			log.Panic(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	l.Info("Shutting down detectr...")
	_ = app.Shutdown()

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
