package main

import (
	"fmt"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	app := fiber.New()

	// Configurar CORS
	app.Use(cors.New())

	// Configure Swagger
	app.Get("/swagger/*", swagger.Handler)

	// Registrar "rooms"
	rooms.Register(app)

	// Initialize Fiber App
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))

}
