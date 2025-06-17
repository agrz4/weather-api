package main

import (
	"errors"
	"go-weather-api/internal/middleware"
	"go-weather-api/service"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

type application struct {
	config         config
	logger         *zap.SugaredLogger
	weatherService service.WeatherService
}

type config struct {
	addr           string
	apiURL         string
	apiKey         string
	redisCfg       redisConfig
	contextTimeout int
}

type redisConfig struct {
	addr string
	pw   string
	db   int
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Add rate limiter middleware
	r.Use(middleware.RateLimiter(rdb))

	r.Route("/health", func(r chi.Router) {
		r.Get("/", app.healthCheckHandler)
	})

	r.Route("/v1", func(r chi.Router) {
		r.Get("/weather", app.weatherHandler)
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("Server has started",
		"addr", app.config.addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
