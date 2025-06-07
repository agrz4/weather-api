package service

import (
	"context"
	"encoding/json"
	"fmt"
	"go-weather-api/internal/model"
	"go-weather-api/internal/repository"
	"go-weather-api/types"
	"go-weather-api/utils"
	"strings"
	"time"

	"go.uber.org/zap"
)

type weatherService struct {
	weatherRepository repository.WeatherRepository
	contextTimeout    time.Duration
	logger            *zap.SugaredLogger
}

type WeatherService interface {
	CreateWeather(ctx context.Context, url *types.Api) (*types.StoreData, error)
	GetWeatherByCity(ctx context.Context, city string) (types.StoreData, error)
}

func NewWeatherService(weatherRepository repository.WeatherRepository, timeout time.Duration, logger *zap.SugaredLogger) WeatherService {
	return &weatherService{
		weatherRepository: weatherRepository,
		contextTimeout:    timeout,
		logger:            logger,
	}
}

func (s *weatherService) CreateWeather(ctx context.Context, url *types.Api) (*types.StoreData, error) {
	var weather *model.Weather

	baseUrl := strings.TrimPrefix(url.Url, "http://")
	apiUrl := fmt.Sprintf("http://%s?key=%s&q=%s", baseUrl, url.ApiKey, url.City)

	body, err := utils.ApiCall(apiUrl)
	if err != nil {
		s.logger.Errorw("Failed to make API call",
			"url", apiUrl,
			"error", err,
		)
		return nil, err
	}

	err = json.Unmarshal(body, &weather)
	if err != nil {
		s.logger.Errorw("Failed to unmarshal API response",
			"response", string(body),
			"error", err,
		)
		return nil, err
	}

	store := &types.StoreData{
		Name:        weather.Location.Name,
		Region:      weather.Location.Region,
		Country:     weather.Location.Country,
		Latitude:    weather.Location.Lat,
		Longitude:   weather.Location.Lon,
		LocalTime:   weather.Location.LocalTime,
		TempC:       weather.Current.TempC,
		TempF:       weather.Current.TempF,
		LastUpdated: weather.Current.LastUpdated,
		Text:        weather.Current.Condition.Text,
		Icon:        weather.Current.Condition.Icon,
		Code:        weather.Current.Condition.Code,
		Uv:          weather.Current.Uv,
	}

	if store.Name == "" {
		s.logger.Errorw("Invalid city name in API response",
			"city", url.City,
			"response", weather,
		)
		return nil, fmt.Errorf("wrong query city name")
	}

	if utils.NormalizeCityName(store.Name) != utils.NormalizeCityName(url.City) {
		s.logger.Errorw("City name mismatch",
			"expected", utils.NormalizeCityName(url.City),
			"actual", utils.NormalizeCityName(store.Name),
		)
		return nil, fmt.Errorf("wrong query city name")
	}

	key := fmt.Sprintf("weather:%s", strings.ToLower(url.City))
	w, err := json.Marshal(store)
	if err != nil {
		s.logger.Errorw("Failed to marshal weather data for caching",
			"data", store,
			"error", err,
		)
		return nil, err
	}
	err = s.weatherRepository.CreateWeather(ctx, key, w)
	if err != nil {
		s.logger.Errorw("Failed to store weather data in cache",
			"key", key,
			"error", err,
		)
		return nil, err
	}
	return store, nil
}

func (s *weatherService) GetWeatherByCity(ctx context.Context, city string) (types.StoreData, error) {
	key := fmt.Sprintf("weather:%s", strings.ToLower(city))
	var data types.StoreData
	cachedWeather, err := s.weatherRepository.GetWeatherByCity(ctx, key)
	if err != nil {
		return types.StoreData{}, err
	}
	err = json.Unmarshal(cachedWeather, &data)
	if err != nil {
		return types.StoreData{}, err
	}
	return data, nil
}
