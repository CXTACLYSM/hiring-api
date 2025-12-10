package clickhouse

import (
	"context"
	"fmt"
	"time"

	clickhouseConfig "github.com/CXTACLYSM/hiring-api/configs/database/analytics/clickhouse"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Connector struct {
	Conn driver.Conn
}

func NewConnector(cfg clickhouseConfig.Config) (*Connector, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.Username,
			Password: cfg.Password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: 5 * time.Second,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("не удалось создать соединение: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		err := conn.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("не удалось подключиться к ClickHouse: %w", err)
	}

	return &Connector{Conn: conn}, nil
}

func (c *Connector) Close() error {
	if c.Conn != nil {
		return c.Conn.Close()
	}
	return nil
}

func (c *Connector) Ping(ctx context.Context) error {
	return c.Conn.Ping(ctx)
}
