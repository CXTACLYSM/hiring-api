package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CXTACLYSM/hiring-api/configs"
	"github.com/CXTACLYSM/hiring-api/internal/di"
	"github.com/CXTACLYSM/hiring-api/internal/routing"
)

func main() {
	cfg, err := configs.Create()
	if err != nil {
		log.Fatalf("Error creating config: %v", err)
	}

	container := di.NewContainer()
	err = container.Init(cfg)
	if err != nil {
		log.Fatalf("Error initializing container: %s", err.Error())
	}
	defer container.PgConnector.Close()

	router := routing.NewRouter()
	err = router.Init(container.Handlers)
	if err != nil {
		log.Fatalf("Error initializing router: %s", err.Error())
	}

	for _, route := range router.GetRoutes() {
		http.Handle(route.String(), route.Handler())
	}

	fmt.Printf("Starting http server on %s\n", cfg.App.SocketStr())
	err = http.ListenAndServe(cfg.App.SocketStr(), nil)
	if err != nil {
		log.Fatalf("error creating http server: %v", err)
	}
}
