package rooms

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golobby/container"
	"github.com/thiagopereiramartinez/scrumpoker-run.api/models/rooms"
	"github.com/thiagopereiramartinez/scrumpoker-run.api/utils"
)

// @Summary Create a new room
// @Tags Rooms
// @Param room body rooms.RoomNewRequest true "Create a new room"
// @Accept json
// @Produce json
// @Success 200 {object} rooms.RoomNewResponse
// @Failure 400 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /rooms [post]
func newRoom(c *fiber.Ctx) error {
	body := new(rooms.RoomNewRequest)
	if err := c.BodyParser(body); err != nil {
		_ = utils.SendError(c, 500, err)
		return nil
	}

	if err := body.Validate(); err != nil {
		_ = utils.SendError(c, 400, err)
		return nil
	}

	db := new(firestore.Client)
	container.Make(&db)

	doc := db.Collection("rooms").NewDoc()
	_, err := doc.Set(context.Background(), map[string]interface{}{
		"name":      body.Name,
		"timestamp": firestore.ServerTimestamp,
	})
	if err != nil {
		_ = utils.SendError(c, 400, err)
		return nil
	}

	return c.JSON(rooms.RoomNewResponse{
		RoomId: doc.ID,
	})
}

// @Summary Join a room
// @Tags Rooms
// @Param room body rooms.RoomJoinRequest true "Join a room"
// @Accept json
// @Produce json
// @Success 200 {object} rooms.RoomJoinResponse
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /rooms/join [post]
func joinRoom(c *fiber.Ctx) error {
	body := new(rooms.RoomJoinRequest)
	if err := c.BodyParser(body); err != nil {
		_ = utils.SendError(c, 500, err)
		return nil
	}

	if err := body.Validate(); err != nil {
		_ = utils.SendError(c, 400, err)
		return nil
	}

	db := new(firestore.Client)
	container.Make(&db)

	roomSnap, err := db.Collection("rooms").Doc(body.RoomId).Get(context.Background())
	if err != nil {
		_ = utils.SendError(c, 404, errors.New("room not found"))
		return nil
	}
	room := new(map[string]interface{})
	err = roomSnap.DataTo(room)
	if err != nil {
		_ = utils.SendError(c, 500, errors.New("unable to retrieve room information"))
		return nil
	}

	player := db.Collection("rooms").Doc(body.RoomId).Collection("players").NewDoc()
	_, err = player.Set(context.Background(), map[string]interface{}{
		"name":      body.Name,
		"timestamp": firestore.ServerTimestamp,
	})
	if err != nil {
		_ = utils.SendError(c, 400, err)
		return nil
	}

	return c.JSON(rooms.RoomJoinResponse{
		RoomId:     body.RoomId,
		RoomName:   fmt.Sprintf("%v", (*room)["name"]),
		PlayerId:   player.ID,
		PlayerName: body.Name,
	})
}

// Registrar endpoints
func Register(app *fiber.App) {
	room := app.Group("/rooms")

	room.Post("", newRoom)
	room.Post("join", joinRoom)
}
