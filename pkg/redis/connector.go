package redis

import (
	"errors"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Db       int

	ConnectTimeout    time.Duration
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
}

type Connector struct {
	AuthPool     *redis.Pool
	ResourcePool *redis.Pool
}

func NewConnector(authCfg, resourceCfg *Config) (*Connector, error) {
	authPool, err := newPool(authCfg)
	if err != nil {
		return nil, fmt.Errorf("error creating auth redis pool: %w", err)
	}

	resourcePool, err := newPool(resourceCfg)
	if err != nil {
		return nil, fmt.Errorf("error creating resource redis pool: %w", err)
	}

	return &Connector{
		AuthPool:     authPool,
		ResourcePool: resourcePool,
	}, nil
}

func newPool(cfg *Config) (*redis.Pool, error) {
	pool := &redis.Pool{
		MaxIdle:     10,
		MaxActive:   50,
		IdleTimeout: time.Minute * 5,
		Dial: func() (redis.Conn, error) {
			dial, err := redis.Dial(
				"tcp",
				fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
				redis.DialUsername(cfg.User),
				redis.DialPassword(cfg.Password),
				redis.DialDatabase(cfg.Db),
				redis.DialConnectTimeout(cfg.ConnectTimeout),
				redis.DialReadTimeout(cfg.ReadTimeout),
				redis.DialWriteTimeout(cfg.WriteTimeout),
			)
			if err != nil {
				return nil, err
			}

			return dial, nil
		},
		TestOnBorrow: func(c redis.Conn, lastUsed time.Time) error {
			if time.Since(lastUsed) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			if err != nil {
				return err
			}

			return nil
		},
	}

	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	if err != nil {
		return nil, fmt.Errorf("error pinging redis: %v", cfg)
	}

	return pool, nil
}

func (c *Connector) Close() error {
	authErr := c.AuthPool.Close()
	resourceErr := c.ResourcePool.Close()

	return errors.Join(authErr, resourceErr)
}
