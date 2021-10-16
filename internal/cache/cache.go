package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var client *redis.Client

// Initialize the client by creating a connection with redis
func Initialize() error {
	client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", viper.GetString("redis_host"), viper.GetString("redis_port")),
	})

	_, err := client.Ping(context.Background()).Result()

	return err
}

func Set(context context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	return client.Set(context, key, value, expiration).Result()
}

func Get(ctx context.Context, key string) (string, error) {
	return client.Get(ctx, key).Result()
}
