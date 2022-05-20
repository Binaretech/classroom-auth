package cache

import (
	"context"
	"fmt"
	"time"

	_ "github.com/Binaretech/classroom-auth/config"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

func init() {
	Initialize()
}

var client *redis.Client

// Initialize the client by creating a connection with redis
func Initialize() error {
	client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", viper.GetString("redis_host"), viper.GetString("redis_port")),
	})

	_, err := client.Ping(context.Background()).Result()

	return err
}

// Set the value of the key in redis with the given expiration time
func Set(context context.Context, key string, value interface{}, expiration time.Duration) (string, error) {
	return client.Set(context, key, value, expiration).Result()
}

// Get the value of the key from redis and return it as a string or an error if the key does not exist
func Get(ctx context.Context, key string) (string, error) {
	return client.Get(ctx, key).Result()
}

func Delete(ctx context.Context, keys ...string) (int64, error) {
	return client.Del(ctx, keys...).Result()
}

func Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	return client.Scan(ctx, cursor, match, count)
}
