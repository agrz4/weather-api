package main

import "go.uber.org/zap"

type application struct {
	config config
	logger *zap.SugaredLogger
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
