package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CXTACLYSM/hiring-api/configs"
	"github.com/CXTACLYSM/hiring-api/internal/handler"
	pgConnector "github.com/CXTACLYSM/hiring-api/pkg/postgres"
)

func main() {
	cfg, err := configs.Create()
	if err != nil {
		log.Fatalf("Failed to load configs: %v", err)
	}

	pgConn, err := pgConnector.NewConnector(cfg.PostgresCluster)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pgConn.Close()

	infoHandler := handler.NewInfoHandler(cfg)
	registerHandler := handler.NewRegisterHandler(pgConn.WritePool)
	authHandler := handler.NewAuthHandler(pgConn.WritePool)

	http.Handle("/", infoHandler)
	http.Handle("/api/register", registerHandler)
	http.Handle("/api/login", authHandler)

	fmt.Printf("Starting http server on %s\n", cfg.App.SocketStr())
	err = http.ListenAndServe(cfg.App.SocketStr(), nil)
	if err != nil {
		log.Fatalf("error creating http server: %v", err)
	}
}
