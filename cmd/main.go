package main

import (
	"fmt"
	"go-weather-api/internal/env"
	"go-weather-api/internal/store"
	"log"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var rdb *redis.Client

func main() {
	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("WEATHER_API_URL", ""),
		apiKey: env.GetString("WEATHER_API_KEY", ""),
		redisCfg: redisConfig{
			addr: env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:   env.GetString("REDIS_PW", ""),
			db:   env.GetInt("REDIS_DB", 0),
		},
		contextTimeout: env.GetInt("TIME_DURATION", 10),
	}
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	rdb = store.NewRedisCache(cfg.redisCfg.addr, cfg.redisCfg.pw, cfg.redisCfg.db)
	logger.Info("redis cache connection established")

	defer rdb.Close()
	app := &application{
		config: cfg,
		logger: logger,
	}

	mux := app.mount()
	if err := app.run(mux); err != nil {
		fmt.Println("err connecting ")
	}
	log.Println(app.run(mux))
}
