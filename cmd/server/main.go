package main

import (
	"os"
	"pollfax/db"
	"pollfax/internal/handlers"
	"pollfax/internal/ingest"
	"pollfax/internal/util"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

func main() {
	util.LoadAppEnv()
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
	e.POST("/bills/sentiments", handlers.CreateBillSentiment)

	e.Logger.Fatal(e.Start(":1323"))
}
