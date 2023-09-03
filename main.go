package main

import (
	// "net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	// "github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"

	"pollfax/database"
	"pollfax/dataingestion"
	// "pollfax/model"
)

func loadAppEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Info().Msg("Error loading .env file")
	}
}

func connectAndMigrateDb() {
	database.ApplyMigrations()
}

func main() {
	loadAppEnv()

	connectAndMigrateDb()

	// dataingestion.ReadData()

	congress := dataingestion.GetLatestCongress()
	dataingestion.GetBills(congress)

	// c := cron.New()
	// c.AddFunc("@daily", func() { dataingestion.Bills() })
	// c.Start()

	r := gin.Default()

	// Apply the CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"*"}
	config.ExposeHeaders = []string{"*"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	r.Use(cors.New(config))

	// r.GET("/ping", func(c *gin.Context) {
	// 	bills, err := model.GetBills()
	// 	if err != nil {
	// 		log.Err(err)
	// 	}
	// 	c.JSON(http.StatusOK, bills)
	// })

	r.Run()
}
