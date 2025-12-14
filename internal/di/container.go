package di

import (
	"fmt"

	"github.com/CXTACLYSM/hiring-api/configs"
	"github.com/CXTACLYSM/hiring-api/internal/handler"
	"github.com/CXTACLYSM/hiring-api/internal/middleware"
	"github.com/CXTACLYSM/hiring-api/internal/routing"
	pgConnector "github.com/CXTACLYSM/hiring-api/pkg/postgres"
)

type Container struct {
	PgConnector *pgConnector.Connector
	Handlers    []routing.HttpRoutable
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) Init(cfg *configs.Config) error {
	pgConn, err := pgConnector.NewConnector(cfg.PostgresCluster)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %v", err)
	}
	c.PgConnector = pgConn

	logging := middleware.Middleware(middleware.Logging)
	authenticate := middleware.Middleware(middleware.Authenticate)

	c.Handlers = make([]routing.HttpRoutable, 0)
	c.Handlers = append(c.Handlers, handler.NewInfoHandler(cfg, []middleware.Middleware{logging, authenticate}))
	c.Handlers = append(c.Handlers, handler.NewRegisterHandler(pgConn, []middleware.Middleware{logging}))
	c.Handlers = append(c.Handlers, handler.NewLoginHandler(pgConn, []middleware.Middleware{logging}))

	return nil
}
