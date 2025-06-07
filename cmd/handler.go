package main

import (
	"context"
	"encoding/json"
	"go-weather-api/types"
	"net/http"
	"time"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	data := map[string]string{
		"status":  "ok",
		"env":     app.config.addr,
		"version": "1.1.0",
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		app.logger.Errorw("Failed to encode health check response", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *application) weatherHandler(w http.ResponseWriter, r *http.Request) {
	// h.logger.Infow("Weather endpoint accessed", "method", r.Method, "url", r.URL.String())

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.config.contextTimeout)*time.Second)
	defer cancel()

	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, `{"status":"error", "message":"city params cannot be empty"}`, http.StatusBadRequest)
		return
	}

	h.logger.Infow("Fetching weather data", "city", city)

	response, err := h.weatherService.GetWeatherByCity(ctx, city)
	// When the Redis key expires, It will be skip err and the data is re-stored in Redis.
	// This prevents continuous API requests with the previously stored key in Redis.
	if err == nil {
		h.logger.Infow("Cached hit", "city", city)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"data":   response,
		})
		return
	}

	// h.logger.Warnw("Cache miss, fetching from API", "city", city)

	apiUrl := &types.Api{
		Url:    h.config.apiURL,
		City:   city,
		ApiKey: h.config.apiKey,
	}

	store, err := h.weatherService.CreateWeather(ctx, apiUrl)
	if err != nil {
		h.logger.Errorw("Failed to fetch weather data from API", "city", city, "error", err)
		http.Error(w, `{"status":"error","message":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	h.logger.Infow("Weather data fetched and cached", "city", city)

	// Return the fetched weather data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"data":   store,
	})
}
