package main

import "github.com/aerosystems/customer-service/internal/handlers"

type Config struct {
	baseHandler *handlers.BaseHandler
}

func NewApp(baseHandler *handlers.BaseHandler) *Config {
	return &Config{
		baseHandler: baseHandler,
	}
}
