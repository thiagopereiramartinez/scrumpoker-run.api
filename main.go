package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

func main() {

	app := fiber.New()

	// Initialize Fiber App
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))

}
