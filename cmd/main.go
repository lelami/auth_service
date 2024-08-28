package main

import (
	_ "authservice/docs"
	"authservice/internal/app"
	"authservice/internal/config"
	"log"
)

// @title           Auth Example API
// @version         0.1.0
// @description     This is a sample auth server with tg auth.

// @contact.name   API Support

// @host      localhost:8000
// @BasePath  /
func main() {
	_, err := config.LoadConfig("../internal/config/config.toml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	app.Run()
}
