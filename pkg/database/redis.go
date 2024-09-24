package db

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// Client variable can used to save key value pairs in redis
var Client *redis.Client

// InitRedis function initializes redis server
func InitRedis() error {
	var err error
	MaxRetries := 5
	RetryDelay := time.Second * 5
	for i := 0; i < MaxRetries; i++ {
		Client = redis.NewClient(&redis.Options{
			Network:  "tcp",
			Addr:     "redis:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		_, err = Client.Ping(ctx).Result()
		if err == nil {
			return nil
		}

		fmt.Printf("Failed to connect to Redis (Attempt %d/%d): %s\n", i+1, MaxRetries, err.Error())
		time.Sleep(RetryDelay)
	}

	// If all attempts fail, return an error
	return fmt.Errorf("failed to connect to Redis after multiple attempts: %s", err.Error())
}

// SetRedis willset a key value in redis server
func SetRedis(key string, value any, expirationTime time.Duration) error {
	fmt.Println("seting value in the reddis !!!!!")
	if err := Client.Set(context.Background(), key, value, expirationTime).Err(); err != nil {
		return err
	}
	return nil
}

// GetRedis will get the value from redis server using key
func GetRedis(key string) (string, error) {
	jsonData, err := Client.Get(context.Background(), key).Result()
	if err != nil {
		return "", err
	}
	return jsonData, nil
}
func DeleteRedis(key string) error {
	fmt.Println("Deleting key from Redis...")
	if err := Client.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}
