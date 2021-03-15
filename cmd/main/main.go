package main

import (
	"cloud.google.com/go/firestore"
	"fmt"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golobby/container"
	_ "github.com/thiagopereiramartinez/scrumpoker-run.api/api"
	"github.com/thiagopereiramartinez/scrumpoker-run.api/internal/controllers/rooms"
	"github.com/thiagopereiramartinez/scrumpoker-run.api/internal/di"
	"log"
	"os"
)

// @title Scrum Poker API
// @version 1.0
func main() {

	// Setup Dependency Injection
	if err := di.SetupDependencies(); err != nil {
		log.Fatalln(err)
	}

	// Create Fiber App
	app := fiber.New()

	// Configure Router
	SetupRouter(app)

	// Initialize Fiber App
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalln(err)
	}

	defer func() {
		// Encerrar conexão com o Firestore
		db := new(firestore.Client)
		container.Make(&db)
		db.Close()
	}()
}

func SetupRouter(app fiber.Router) {

	// Setup CORS
	app.Use(cors.New())

	// Setup Swagger
	app.Get("/swagger", func(ctx *fiber.Ctx) error {
		return ctx.Redirect("/swagger/index.html")
	})
	app.Get("/swagger/*", swagger.Handler)

	// Register "rooms"
	rooms.Register(app)

}
