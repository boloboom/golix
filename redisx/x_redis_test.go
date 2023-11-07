package redisx

import (
	"context"
	"log"
	"testing"
	"time"
)

var testSetting = &RedisSetting{
	Host:        "localhost:6379",
	Username:    "default",
	Password:    "123456",
	DB:          0,
	MaxIdle:     30,
	MaxActive:   30,
	IdleTimeout: 200,
	Prefix:      "xiaoqucloud",
}

var testKey = CacheKey("test_key")

func TestRedis(t *testing.T) {
	Setup(testSetting)
	Cache().Client.Set(context.Background(), testKey.Key(), "test_value", time.Minute)
	result, err := Cache().Client.Get(context.Background(), testKey.Key()).Result()
	if err != nil {
		t.Error(err)
	}
	log.Println(result)
}
