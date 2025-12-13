package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CXTACLYSM/hiring-api/configs/database/persistence/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Connector struct {
	ReadPool  *pgxpool.Pool
	WritePool *pgxpool.Pool
}

func NewConnector(cfg *postgres.ClusterConfig) (*Connector, error) {
	readDSN, err := cfg.DSN(postgres.ReadOperation)
	if err != nil {
		return nil, err
	}
	readPool, err := getPool(readDSN)
	if err != nil {
		return nil, err
	}

	writeDSN, err := cfg.DSN(postgres.WriteOperation)
	if err != nil {
		return nil, err
	}
	writePool, err := getPool(writeDSN)
	if err != nil {
		return nil, err
	}

	return &Connector{
		ReadPool:  readPool,
		WritePool: writePool,
	}, nil
}

func getPool(dsn string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("error parsing DSN: %w", err)
	}

	poolCfg.MaxConns = 25
	poolCfg.MinConns = 5
	poolCfg.MaxConnLifetime = time.Hour
	poolCfg.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("cannot create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("cannot connect to PostgreSQL: %w", err)
	}
	fmt.Printf("successfully connected to PostgreSQL at %s", dsn)

	return pool, nil
}

func (c *Connector) Close() {
	c.ReadPool.Close()
	c.WritePool.Close()
}

func (c *Connector) PingByOperation(operation uint8, ctx context.Context) error {
	switch true {
	case operation == postgres.ReadOperation && c.ReadPool != nil:
		err := c.ReadPool.Ping(ctx)
		if err != nil {
			return err
		}
		return nil
	case operation == postgres.WriteOperation && c.WritePool != nil:
		err := c.WritePool.Ping(ctx)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("operation not supported or corresponding pool is nil")
	}
}
