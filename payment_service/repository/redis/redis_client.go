package redis

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func InitRedisClient(addr string, password string, db int) {
	log.Printf("Initializing Redis client with address: %s", addr)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	pong, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Printf("ERROR: Failed to connect to Redis at %s: %v", addr, err)
	} else {
		log.Printf("SUCCESS: Redis connection established at %s. Response: %s", addr, pong)
	}

	testKey := "test:connection"
	testValue := "redis_is_working"

	err = RedisClient.Set(Ctx, testKey, testValue, 1*time.Minute).Err()
	if err != nil {
		log.Printf("ERROR: Failed to set test key in Redis: %v", err)
	} else {
		log.Printf("Redis SET operation successful")

		val, err := RedisClient.Get(Ctx, testKey).Result()
		if err != nil {
			log.Printf("ERROR: Failed to get test key from Redis: %v", err)
		} else if val == testValue {
			log.Printf("Redis GET operation successful. Redis is fully operational.")
		} else {
			log.Printf("WARNING: Redis GET returned unexpected value: %s", val)
		}
	}
}
