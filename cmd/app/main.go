package main

import (
	"Norvista/api"
	"Norvista/internal/config"
	"Norvista/internal/database"
	"fmt"

	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Info().Msg("Starting Norvista project")
	connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable",
		config.Envs.DBUser,
		config.Envs.DBPassword,
		config.Envs.DBAddress,
		config.Envs.DBName,
	)
	sqlStorage, err := database.NewPostgresStorage(connStr)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	db, err := sqlStorage.InitializeDatabase()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}

	store := api.NewStore(db)
	apiServer := api.NewAPIServer(":8080", store)
	log.Info().Msg("Starting API server on port 8080")
	apiServer.Serve()
}
