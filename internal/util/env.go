package util

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func LoadAppEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}
}
