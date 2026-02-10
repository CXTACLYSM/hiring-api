package di

import (
	"fmt"

	"github.com/CXTACLYSM/hiring-api/configs"
	"github.com/CXTACLYSM/hiring-api/internal/auth/tokens"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/commands/createOne"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/queries/findOne"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/services"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/infrastructure/handlers"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/infrastructure/middlewares"
	pgConnector "github.com/CXTACLYSM/hiring-api/pkg/postgres"
)

type Infrastructure struct {
	PgConnector *pgConnector.Connector
}

type Queries struct {
	findOneUser   *findOne.Handler
	createOneUser *createOne.Handler
}

type Services struct {
	authService *services.AuthService
}

type Handlers struct {
	Register *handlers.RegisterHandler
	Login    *handlers.LoginHandler
	Me       *handlers.MeHandler
}

type Middlewares struct {
	Authenticate *middlewares.Authenticate
	Json         *middlewares.ContentType
}

type Container struct {
	Infrastructure *Infrastructure
	Queries        *Queries
	Services       *Services
	Middlewares    *Middlewares
	Handlers       *Handlers
}

func NewContainer() *Container {
	return &Container{
		Middlewares: &Middlewares{},
		Handlers:    &Handlers{},
	}
}

func (c *Container) Init(cfg *configs.Config) error {
	if err := c.initInfrastructure(cfg); err != nil {
		return err
	}
	if err := c.initQueries(); err != nil {
		return err
	}
	if err := c.initMiddlewares(cfg); err != nil {
		return err
	}
	if err := c.initServices(cfg); err != nil {
		return err
	}
	if err := c.initHandlers(); err != nil {
		return err
	}

	return nil
}

func (c *Container) initInfrastructure(cfg *configs.Config) error {
	pgConn, err := pgConnector.NewConnector(cfg.PostgresCluster)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %v", err)
	}

	c.Infrastructure = &Infrastructure{
		PgConnector: pgConn,
	}

	return nil
}

func (c *Container) initQueries() error {
	c.Queries = &Queries{
		findOneUser:   findOne.NewFindOneUserQueryHandler(c.Infrastructure.PgConnector),
		createOneUser: createOne.NewCreateOneUserQueryHandler(c.Infrastructure.PgConnector),
	}

	return nil
}

func (c *Container) initMiddlewares(cfg *configs.Config) error {
	c.Middlewares = &Middlewares{
		Authenticate: middlewares.NewAuthenticate([]byte(cfg.App.JwtSecret)),
		Json:         middlewares.NewContentType("application/json"),
	}

	return nil
}

func (c *Container) initServices(cfg *configs.Config) error {
	tokenGenerator := tokens.NewJwtTokenGenerator(cfg.App.JwtSecret)

	c.Services = &Services{
		authService: services.NewAuthService(c.Queries.findOneUser, c.Queries.createOneUser, tokenGenerator),
	}

	return nil
}

func (c *Container) initHandlers() error {
	c.Handlers.Register = handlers.NewRegisterHandler(c.Services.authService)
	c.Handlers.Login = handlers.NewLoginHandler(c.Services.authService)

	return nil
}
