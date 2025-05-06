package main

import (
	"Solution/config"
	"Solution/internal/app"
	"go.uber.org/zap"
	"log"
)

func main() {
	cfg, err := config.New("config.json")

	if err != nil {
		log.Fatalf("can not get application config: %s", err)
	}

	var logger *zap.Logger

	logger, err = zap.NewProduction()

	if err != nil {
		log.Fatalf("can not initialize logger: %s", err)
	}

	app.Run(logger, cfg)
}
