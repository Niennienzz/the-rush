package repository

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"

	"the-rush/constant"
)

type Repository struct {
	*redis.Client
	maxLimit                       int
	playerRecords                  string
	playerInterimPrefix            string
	playerByNamePrefix             string
	playerByCreatedBy              string
	playerByLongestRush            string
	playerByTotalRushingTouchdowns string
	playerByTotalRushingYards      string
}

func NewLocal() *Repository {
	client := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_URL", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return &Repository{
		Client:                         client,
		maxLimit:                       constant.MaxPageLimit,
		playerRecords:                  constant.RedisKeyPlayerRecords.String(),
		playerInterimPrefix:            constant.RedisIndexKeyPlayerInterimPrefix.String(),
		playerByNamePrefix:             constant.RedisIndexKeyPlayerByNamePrefix.String(),
		playerByCreatedBy:              constant.RedisIndexKeyPlayerByCreatedAt.String(),
		playerByLongestRush:            constant.RedisIndexKeyPlayerByLongestRush.String(),
		playerByTotalRushingTouchdowns: constant.RedisIndexKeyPlayerByTotalRushingTouchdowns.String(),
		playerByTotalRushingYards:      constant.RedisIndexKeyPlayerByTotalRushingYards.String(),
	}
}

func (x *Repository) Player() PlayerRepository {
	return &playerRepository{x}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
