package main

import (
	"authservice/internal/app"
	"authservice/internal/config"
	"log"
)

func main() {
	_, err := config.LoadConfig("../internal/config/config.toml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	app.Run()
}
