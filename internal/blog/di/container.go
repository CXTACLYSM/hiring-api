package di

import (
	"fmt"

	"github.com/CXTACLYSM/hiring-api/configs/blog"
	"github.com/CXTACLYSM/hiring-api/configs/blog/app"
	createOnePost "github.com/CXTACLYSM/hiring-api/internal/blog/post/application/commands/createOne"
	deleteOnePost "github.com/CXTACLYSM/hiring-api/internal/blog/post/application/commands/deleteOne"
	updateOnePost "github.com/CXTACLYSM/hiring-api/internal/blog/post/application/commands/updateOne"
	findListPost "github.com/CXTACLYSM/hiring-api/internal/blog/post/application/queries/findList"
	findOnePost "github.com/CXTACLYSM/hiring-api/internal/blog/post/application/queries/findOne"
	postServices "github.com/CXTACLYSM/hiring-api/internal/blog/post/application/services"
	createOnePostHandler "github.com/CXTACLYSM/hiring-api/internal/blog/post/infrastructure/commands/createOne"
	deleteOnePostHandler "github.com/CXTACLYSM/hiring-api/internal/blog/post/infrastructure/commands/deleteOne"
	updateOnePostHandler "github.com/CXTACLYSM/hiring-api/internal/blog/post/infrastructure/commands/updateOne"
	"github.com/CXTACLYSM/hiring-api/internal/blog/post/infrastructure/handlers"
	findListPostHandler "github.com/CXTACLYSM/hiring-api/internal/blog/post/infrastructure/queries/findList"
	findOnePostHandler "github.com/CXTACLYSM/hiring-api/internal/blog/post/infrastructure/queries/findOne"
	pb "github.com/CXTACLYSM/hiring-api/pkg/grpc/auth/v1"
	pgConnector "github.com/CXTACLYSM/hiring-api/pkg/postgres"
	redisConnector "github.com/CXTACLYSM/hiring-api/pkg/redis"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/cache"
	pkgHandlers "github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/handlers"
	pkgMiddlewares "github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/middlewares"
	"github.com/CXTACLYSM/hiring-api/pkg/shared/infrastructure/validation"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Infrastructure struct {
	PgConnector    *pgConnector.Connector
	RedisConnector *redisConnector.Connector
	AuthGrpcConn   *grpc.ClientConn
	Logger         *zap.Logger
}

type Queries struct {
	FindOnePost  findOnePost.Handler
	FindListPost findListPost.Handler
}

type Commands struct {
	CreateOnePost createOnePost.Handler
	UpdateOnePost updateOnePost.Handler
	DeleteOnePost deleteOnePost.Handler
}

type CacheEntities struct {
	UserCache *cache.UserCache
}

type Services struct {
	PostService *postServices.PostService
}

type Handlers struct {
	Info          *pkgHandlers.InfoHandler
	FindListPost  *handlers.FindListPostHandler
	FindOnePost   *handlers.FindOnePostHandler
	CreateOnePost *handlers.CreateOnePostHandler
	UpdateOnePost *handlers.UpdateOnePostHandler
	DeleteOnePost *handlers.DeleteOnePostHandler
}

type Middlewares struct {
	Authenticate *pkgMiddlewares.Authenticate
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
	if err := c.initValidator(); err != nil {
		return err
	}
	if err := c.initCacheEntities(); err != nil {
		return err
	}

	if err := c.initMiddlewares(); err != nil {
		return err
	}
	if err := c.initServices(); err != nil {
		return err
	}
	if err := c.initHandlers(cfg); err != nil {
		return err
	}

	return nil
}

func (c *Container) initInfrastructure(cfg *configs.Config) error {
	logger, err := createLogger(cfg.App.Environment)
	if err != nil {
		return fmt.Errorf("error creating zap logger: %w", err)
	}

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
		return fmt.Errorf("faied to connect to redis: %w", err)
	}

	conn, err := grpc.NewClient(cfg.Auth.GrpcSocketStr(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to auth grpc: %w", err)
	}

	c.Infrastructure = &Infrastructure{
		PgConnector:    pgConn,
		RedisConnector: redisConn,
		AuthGrpcConn:   conn,
		Logger:         logger,
	}

	return nil
}

func (c *Container) initQueries() error {
	c.Queries = &Queries{
		FindOnePost:  findOnePostHandler.NewFindOnePostHandler(c.Infrastructure.PgConnector.ReadPool),
		FindListPost: findListPostHandler.NewFindListPostHandler(c.Infrastructure.PgConnector.ReadPool),
	}

	return nil
}

func (c *Container) initCommands() error {
	c.Commands = &Commands{
		CreateOnePost: createOnePostHandler.NewCreateOnePostHandler(c.Infrastructure.PgConnector.WritePool),
		UpdateOnePost: updateOnePostHandler.NewUpdateOnePostHandler(c.Infrastructure.PgConnector.WritePool),
		DeleteOnePost: deleteOnePostHandler.NewDeleteOnePostHandler(c.Infrastructure.PgConnector.WritePool),
	}

	return nil
}

func (c *Container) initCacheEntities() error {
	c.CacheEntities = &CacheEntities{
		UserCache: cache.NewUserCache(c.Infrastructure.RedisConnector.AuthPool),
	}

	return nil
}

func (c *Container) initMiddlewares() error {
	c.Middlewares = &Middlewares{
		Authenticate: pkgMiddlewares.NewAuthenticate(c.CacheEntities.UserCache, pb.NewAuthServiceClient(c.Infrastructure.AuthGrpcConn)),
		Json:         pkgMiddlewares.NewContentType("application/json"),
	}
	return nil
}

func (c *Container) initValidator() error {
	c.Validator = validation.NewValidator(validator.New())

	return nil
}

func (c *Container) initServices() error {
	c.Services = &Services{
		PostService: postServices.NewPostService(
			c.Validator,
			c.Queries.FindOnePost,
			c.Queries.FindListPost,
			c.Commands.CreateOnePost,
			c.Commands.UpdateOnePost,
			c.Commands.DeleteOnePost,
			c.Infrastructure.Logger,
		),
	}

	return nil
}

func (c *Container) initHandlers(cfg *configs.Config) error {
	c.Handlers = &Handlers{
		Info:          pkgHandlers.NewInfoHandler(cfg.App.Version),
		FindListPost:  handlers.NewFindListPostHandler(c.Services.PostService),
		FindOnePost:   handlers.NewFindOnePostHandler(c.Services.PostService),
		CreateOnePost: handlers.NewCreateOnePostHandler(c.Services.PostService),
		UpdateOnePost: handlers.NewUpdateOnePostHandler(c.Services.PostService),
		DeleteOnePost: handlers.NewDeleteOnePostHandler(c.Services.PostService),
	}

	return nil
}

func createLogger(env string) (*zap.Logger, error) {
	if env == app.Production {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}
