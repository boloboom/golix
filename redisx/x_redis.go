package redisx

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisSetting struct {
	Host     string
	Username string
	Password string
	DB       int

	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration

	Prefix string
}

type RedisCache struct {
	Client *redis.Client
	prefix string
}

var (
	rc   RedisCache
	once sync.Once
)

// Setup initializes the Redis client with the given RedisSetting.
//
// It takes a pointer to a RedisSetting struct as a parameter.
// It does not return anything.
func Setup(setting *RedisSetting) {
	once.Do(func() {
		rc.prefix = setting.Prefix
		rc.Client = redis.NewClient(&redis.Options{
			Addr:            setting.Host,
			PoolSize:        setting.MaxActive,
			MaxIdleConns:    setting.MaxIdle,
			ConnMaxIdleTime: setting.IdleTimeout,
			Username:        setting.Username,
			Password:        setting.Password,
			DB:              setting.DB,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		_, err := rc.Client.Ping(ctx).Result()
		if err != nil {
			log.Println(err.Error())
			log.Fatal("redis init failed, please check redis config")
		}
	})
}

// Cache returns a pointer to a RedisCache instance.
//
// No parameters.
// Returns a pointer to a RedisCache instance.
func Cache() *RedisCache {
	return &rc
}
