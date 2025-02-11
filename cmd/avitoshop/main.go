package main

import (
	"log"

	app "github.com/RomanAgaltsev/avito-shop/internal/app/avitoshop"
	"github.com/RomanAgaltsev/avito-shop/internal/config"
)

func main() {
	// Get application cofiguration
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("config building: %s", err)
	}

	// Run the application
	err = app.Run(cfg)
	if err != nil {
		log.Fatalf("failed to run application : %s", err.Error())
	}
}
