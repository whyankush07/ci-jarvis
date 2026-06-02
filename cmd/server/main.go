package main

import (
	"ci-jarvis/api"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", api.HealthCheck)
	app.Post("/webhook", api.WebhookHandler)

	addr := ":8080"
	log.Printf("Server has been configured for %s\n", addr)
	if err := app.Listen(addr); err != nil {
		fmt.Println("Error in starting the server!!", err)
	}
}
