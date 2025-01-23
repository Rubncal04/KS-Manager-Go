package db

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Rubncal04/ksmanager/config"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	re *redis.Client
}

func StartCacheDb(variables *config.EnvVariables) *RedisRepo {
	var ctx = context.Background()
	db := 0
	if variables.REDIS_DATABASES != "" {
		db, _ = strconv.Atoi(variables.REDIS_DATABASES)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     variables.REDIS_ADDRESS + ":" + variables.REDIS_PORT,
		Password: variables.REDIS_PASSWORD,
		DB:       db,
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis can't be connected: %v", err)
	}

	fmt.Printf("Awesome! You're connected to Redis 🎊🥳: %s\n", pong)

	return &RedisRepo{
		re: rdb,
	}
}

func (r *RedisRepo) GetCache(cacheKey string) (string, error) {
	ctx := context.Background()

	cachedData, err := r.re.Get(ctx, cacheKey).Result()
	if err != nil {
		log.Printf("Error getting data from cache: %v", err)
		return "", err
	}

	return cachedData, nil
}

func (r *RedisRepo) SetCache(cacheKey string, recorder string) (string, error) {
	ctx := context.Background()
	err := r.re.Set(ctx, cacheKey, recorder, time.Hour).Err()

	if err != nil {
		log.Printf("Error setting cache: %v", err)
		return "", err
	}

	return "Data stored in cache", nil
}

func (r *RedisRepo) DeleteCache(cacheKey string) (int64, error) {
	ctx := context.Background()

	deleted, err := r.re.Del(ctx, cacheKey).Result()
	if err != nil {
		log.Printf("Error deleting cache: %v", err)
		return 0, err
	}
	return deleted, nil
}
