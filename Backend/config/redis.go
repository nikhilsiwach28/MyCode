package config

import (
	"log"
	"strconv"
)

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

func NewRedisConfig() RedisConfig {
	db, err := strconv.Atoi(GetEnvWithDefault("REDIS_DB", "0"))
	if err != nil {
		log.Fatal(err)
	}
	return RedisConfig{
		Address:  GetEnvWithDefault("REDIS_ADDRESS", "localhost:6379"),
		Password: GetEnvWithDefault("REDIS_PASSWORD", ""),
		DB:       db,
	}
}
