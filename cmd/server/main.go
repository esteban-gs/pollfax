package main

import (
	"os"
	"pollfax/db"
	"pollfax/internal/handlers"
	"pollfax/internal/ingest"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func loadAppEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Info().Msg("Error loading .env file")
	}
}

func main() {
	loadAppEnv()

	db.ApplyMigrations()

	// RUN DATA PIPELINE NIGHTLY
	c := cron.New()
	c.AddFunc("@daily", func() { ingest.Run() })
	c.Start()

	e := echo.New()
	e.Use(middleware.CORS())

	logger := zerolog.New(os.Stdout)

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")

			return nil
		},
	}))

	e.GET("/bills", handlers.GetAll)
	e.File("/bills.json", "public/bills.json")

	e.Logger.Fatal(e.Start(":1323"))
}
