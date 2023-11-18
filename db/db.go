package db

import (
	"database/sql"
	"fmt"

	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func ConnectionString() string {
	host := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	databaseName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, username, password, databaseName, port, sslmode)
}

func ApplyMigrations() {
	databaseName := os.Getenv("DB_NAME")
	dsn := ConnectionString()
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("Error applying migrations")
		panic(err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		databaseName,
		driver)
	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
	log.Info().Msg("Successfully applied migrations")
}

func Instance() *sqlx.DB {
	dsn := ConnectionString()
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
		panic(err)
	}
	return db
}
