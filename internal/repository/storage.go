package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type WeatherRepository interface {
	CreateWeather(ctx context.Context, key string, weather []byte) error
	GetWeatherByCity(ctx context.Context, key string) ([]byte, error)
}

func NewWeatherRepo(client *redis.Client, expiry time.Duration) WeatherRepository {
	return &weatherRepository{
		client: client,
		expiry: expiry,
	}
}
