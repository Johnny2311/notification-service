package main

import (
	_ "github.com/joho/godotenv/autoload" // Load variables from .env file.
	"notification-service/cmd/notifications/server"
)

func main() {
	server.Run()
}
