package main

import (
	"authservice/config"
	"authservice/internal/app"
	"log"
)

func main() {
	_, err := config.LoadConfig("../config/config.toml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	app.Run()
}
