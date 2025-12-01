package main

import (
	"log"

	"any-x/internal/app"
	"any-x/internal/config"
)

func main() {
	cfg := config.DefaultLocal()

	application := app.New(cfg)

	if err := application.Run(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
