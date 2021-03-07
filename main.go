package main

import (
	"fmt"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/thiagopereiramartinez/scrumpoker-run.api/di"
	_ "github.com/thiagopereiramartinez/scrumpoker-run.api/docs"
	"github.com/thiagopereiramartinez/scrumpoker-run.api/rooms"
	"log"
	"os"
)

// @title Scrum Poker API
// @version 1.0
// @contact.name Thiago P. Martinez
// @contact.email thiago.pereira.ti@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
func main() {

	// Setup Dependency Injection
	if err := di.SetupDependencies(); err != nil {
		log.Fatalln(err)
	}

	// Create Fiber App
	app := fiber.New()

	// Setup CORS
	app.Use(cors.New())

	// Setup Swagger
	app.Get("/swagger/*", swagger.Handler)

	// Register "rooms"
	rooms.Register(app)

	// Initialize Fiber App
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalln(err)
	}

}
