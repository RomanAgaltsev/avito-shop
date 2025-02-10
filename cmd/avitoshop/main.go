package main

import (
	"log"
	
	"github.com/RomanAgaltsev/avito-shop/internal/config"
)

func main() {
	// Get application cofiguration
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("config building: %s", err)
	}
	// Run application
	app.Run(cfg)
}
