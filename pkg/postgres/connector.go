package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/CXTACLYSM/hiring-api/configs/database/persistence/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Connector struct {
	Pool *pgxpool.Pool
}

func NewConnector(cfg postgres.Config) (*Connector, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("error parsing DSN: %w", err)
	}

	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("cannot connect to PostgreSQL: %w", err)
	}

	return &Connector{Pool: pool}, nil
}

func (c *Connector) Close() {
	if c.Pool != nil {
		c.Pool.Close()
	}
}

func (c *Connector) Ping(ctx context.Context) error {
	return c.Pool.Ping(ctx)
}
