package main

import (
	"Norvista/api"
	"Norvista/internal/config"
	"Norvista/internal/database"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Start Norvista project")
	connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable",
		config.Envs.DBUser,
		config.Envs.DBPassword,
		config.Envs.DBAddress,
		config.Envs.DBName,
	)
	sqlStorage, err := database.NewPostgresStorage(connStr)
	if err != nil {
		log.Fatal(err)
	}

	 db, err := sqlStorage.InitializeDatabase();
	 if err != nil {
		log.Fatal(err)
	}

	store := api.NewStore(db)
	apiServer := api.NewAPIServer(":8080", store)
	apiServer.Serve()
}
