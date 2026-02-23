package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	entities "github.com/CXTACLYSM/hiring-api/pkg/shared/domain"
	"github.com/gomodule/redigo/redis"
)

const (
	userCacheKeyPrefix = "auth:token:"
)

const (
	userCacheTTL = 5 * 60
)

type UserCacheKey string

type UserCache struct {
	pool *redis.Pool
}

func NewUserCache(pool *redis.Pool) *UserCache {
	return &UserCache{
		pool: pool,
	}
}

func (c *UserCache) Key(token string) UserCacheKey {
	return UserCacheKey(fmt.Sprintf("%s:%s", userCacheKeyPrefix, token))
}

func (c *UserCache) Get(key UserCacheKey) (*entities.User, bool) {
	conn := c.pool.Get()
	defer conn.Close()

	data, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		if !errors.Is(err, redis.ErrNil) {
			log.Printf("redis GET error for key=%s: %v", key, err)
		}
		return nil, false
	}

	var user entities.User
	if err = json.Unmarshal(data, &user); err != nil {
		log.Printf("user unmarshal error form redis for key=%s: %v", key, err)
		conn.Do("DEL", key)
		return nil, false
	}

	return &user, true
}

func (c *UserCache) Set(key UserCacheKey, user *entities.User) error {
	if user == nil {
		return fmt.Errorf("error setting nil user to redis for key=%s", key)
	}
	conn := c.pool.Get()
	defer conn.Close()

	data, err := json.Marshal(user)
	if err != nil {
		log.Printf("redis cache marshal error for key=%s: %v", key, err)
		return fmt.Errorf("redis cache marshal error for key=%s: %w", key, err)
	}

	_, err = conn.Do("SET", key, data, "EX", userCacheTTL)
	if err != nil {
		return fmt.Errorf("error setting cache key=%s: %v", key, err)
	}

	return nil
}
