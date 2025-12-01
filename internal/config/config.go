package config

import (
	"time"
)

type Config struct {
	Env         string
	StoragePath string
	HTTPServer  HTTPServer
}

type HTTPServer struct {
	Address     string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

func DefaultLocal() *Config {
	return &Config{
		Env:         "local",
		StoragePath: "./storage.db",
		HTTPServer: HTTPServer{
			Address:     ":8080",
			Timeout:     4 * time.Second,
			IdleTimeout: 30 * time.Second,
		},
	}
}
