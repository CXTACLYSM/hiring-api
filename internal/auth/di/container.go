package di

import (
	"fmt"

	"github.com/CXTACLYSM/hiring-api/configs/auth"
	"github.com/CXTACLYSM/hiring-api/internal/auth/shared/infrastructure/middlewares"
	"github.com/CXTACLYSM/hiring-api/internal/auth/tokens"
	createOneUser "github.com/CXTACLYSM/hiring-api/internal/auth/user/application/commands/createOne"
	findOneUser "github.com/CXTACLYSM/hiring-api/internal/auth/user/application/queries/findOne"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/application/services"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/infrastructure/commands/createOne"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/infrastructure/handlers"
	"github.com/CXTACLYSM/hiring-api/internal/auth/user/infrastructure/queries/findOne"
	"github.com/CXTACLYSM/hiring-api/pkg/kafka"
	pgConnector "github.com/CXTACLYSM/hiring-api/pkg/postgres"
	redisConnector "github.com/CXTACLYSM/hiring-api/pkg/redis"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/cache"
	pkgHandlers "github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/handlers"
	pkgMiddlewares "github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/middlewares"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/validation"
	"github.com/IBM/sarama"
	"github.com/go-playground/validator/v10"
)

type Infrastructure struct {
	PgConnector    *pgConnector.Connector
	RedisConnector *redisConnector.Connector
	Kafka          *Kafka
}

type Kafka struct {
	SyncProducer sarama.SyncProducer
	Publisher    *kafka.Publisher
}

type Queries struct {
	FindOneUser findOneUser.Handler
}

type Commands struct {
	CreateOneUser createOneUser.Handler
}

type CacheEntities struct {
	UserCache *cache.UserCache
}

type Services struct {
	authService *services.AuthService
	userService *services.UserService
}

type Handlers struct {
	Info     *pkgHandlers.InfoHandler
	Register *handlers.RegisterHandler
	Login    *handlers.LoginHandler
	Me       *handlers.MeHandler
}

type Middlewares struct {
	Authenticate *middlewares.Authenticate
	Json         *pkgMiddlewares.ContentType
}

type Container struct {
	Infrastructure *Infrastructure
	Queries        *Queries
	Commands       *Commands
	CacheEntities  *CacheEntities
	Services       *Services
	Middlewares    *Middlewares
	Handlers       *Handlers
	Validator      *validation.Validator
}

func (c *Container) Init(cfg *configs.Config) error {
	if err := c.initInfrastructure(cfg); err != nil {
		return err
	}
	if err := c.initQueries(); err != nil {
		return err
	}
	if err := c.initCommands(); err != nil {
		return err
	}
	if err := c.initCacheEntities(); err != nil {
		return err
	}
	if err := c.initValidator(); err != nil {
		return err
	}
	if err := c.initMiddlewares(cfg); err != nil {
		return err
	}
	if err := c.initServices(cfg); err != nil {
		return err
	}
	if err := c.initHandlers(cfg); err != nil {
		return err
	}

	return nil
}

func (c *Container) initInfrastructure(cfg *configs.Config) error {
	readDSN, err := cfg.PostgresCluster.DSN(pgConnector.ReadOperation)
	if err != nil {
		return fmt.Errorf("error initializing infra read pgx pool: %w", err)
	}
	writeDSN, err := cfg.PostgresCluster.DSN(pgConnector.WriteOperation)
	if err != nil {
		return fmt.Errorf("error initializing infra write pgx pool: %w", err)
	}
	pgConn, err := pgConnector.NewConnector(readDSN, writeDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}

	authCfg, resourceCfg := cfg.Redis.ConnectorConfigs()
	redisConn, err := redisConnector.NewConnector(authCfg, resourceCfg)
	if err != nil {
		return fmt.Errorf("failed to connecto to redis: %w", err)
	}

	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	producer, err := sarama.NewSyncProducer(cfg.Kafka.Brokers(), config)
	if err != nil {
		return err
	}
	publisher := kafka.NewPublisher(producer)

	c.Infrastructure = &Infrastructure{
		PgConnector:    pgConn,
		RedisConnector: redisConn,
		Kafka: &Kafka{
			SyncProducer: producer,
			Publisher:    publisher,
		},
	}

	return nil
}

func (c *Container) initQueries() error {
	c.Queries = &Queries{
		FindOneUser: findOne.NewFindOneUserQueryHandler(c.Infrastructure.PgConnector.ReadPool),
	}

	return nil
}

func (c *Container) initCommands() error {
	c.Commands = &Commands{
		CreateOneUser: createOne.NewCreateOneUserQueryHandler(c.Infrastructure.PgConnector.WritePool),
	}

	return nil
}

func (c *Container) initCacheEntities() error {
	c.CacheEntities = &CacheEntities{
		UserCache: cache.NewUserCache(c.Infrastructure.RedisConnector.AuthPool),
	}

	return nil
}

func (c *Container) initMiddlewares(cfg *configs.Config) error {
	c.Middlewares = &Middlewares{
		Authenticate: middlewares.NewAuthenticate([]byte(cfg.App.JwtSecret), c.CacheEntities.UserCache, c.Queries.FindOneUser),
		Json:         pkgMiddlewares.NewContentType("application/json"),
	}

	return nil
}

func (c *Container) initValidator() error {
	c.Validator = validation.NewValidator(validator.New())

	return nil
}

func (c *Container) initServices(cfg *configs.Config) error {
	tokenGenerator := tokens.NewJwtTokenGenerator(cfg.App.JwtSecret)

	c.Services = &Services{
		authService: services.NewAuthService(c.Validator, c.Queries.FindOneUser, c.Commands.CreateOneUser, tokenGenerator, c.Infrastructure.Kafka.Publisher),
		userService: services.NewUserService(c.Queries.FindOneUser, c.Commands.CreateOneUser),
	}

	return nil
}

func (c *Container) initHandlers(cfg *configs.Config) error {
	c.Handlers = &Handlers{
		Info:     pkgHandlers.NewInfoHandler(cfg.App.Version),
		Register: handlers.NewRegisterHandler(c.Services.authService),
		Login:    handlers.NewLoginHandler(c.Services.authService),
		Me:       handlers.NewMeHandler(c.Services.userService),
	}

	return nil
}
