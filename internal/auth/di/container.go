package di

import (
	"fmt"

	"github.com/CXTACLYSM/hiring-api/configs"
	"github.com/CXTACLYSM/hiring-api/internal/auth/handlers"
	"github.com/CXTACLYSM/hiring-api/internal/auth/middlewares"
	"github.com/CXTACLYSM/hiring-api/internal/auth/queries/user/createOne"
	"github.com/CXTACLYSM/hiring-api/internal/auth/queries/user/findOne"
	"github.com/CXTACLYSM/hiring-api/internal/auth/routing"
	"github.com/CXTACLYSM/hiring-api/internal/auth/services"
	"github.com/CXTACLYSM/hiring-api/internal/auth/tokens"
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

	// middlewares
	logging := middlewares.Middleware(middlewares.Logging)
	//authenticate := middleware.Middleware(middleware.Authenticate)

	// token generator
	tokenGenerator := tokens.NewJwtTokenGenerator(cfg.App.JwtSecret)

	// query handlers
	findOneUserQueryHandler := findOne.NewFindOneUserQueryHandler(c.PgConnector)
	createOneUserQueryHandler := createOne.NewCreateOneUserQueryHandler(c.PgConnector)

	// services
	authService := services.NewAuthService(findOneUserQueryHandler, createOneUserQueryHandler, services.TokenGenerator(tokenGenerator))

	// handlers
	c.Handlers = make([]routing.HttpRoutable, 0)
	c.Handlers = append(c.Handlers, handlers.NewInfoHandler(cfg, logging))
	c.Handlers = append(c.Handlers, handlers.NewRegisterHandler(authService, logging))
	c.Handlers = append(c.Handlers, handlers.NewLoginHandler(authService, logging))

	return nil
}
