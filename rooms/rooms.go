package rooms

import "github.com/gofiber/fiber/v2"

// @Summary Create a new room
// @Tags Rooms
// @Accept json
// @Accept json
// @Router /rooms [post]
func newRoom(c *fiber.Ctx) error {
	return c.SendString("Deu certo")
}

func Register(app *fiber.App) {
	room := app.Group("/rooms")

	room.Post("", newRoom)
}
