package rooms

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golobby/container"
	"github.com/thiagopereiramartinez/scrumpoker-run.api/models/players"
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
// @Param body body rooms.RoomJoinRequest true "Join a room"
// @Param id path string true "Room Id"
// @Accept json
// @Produce json
// @Success 200 {object} rooms.RoomJoinResponse
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /rooms/{id}/join [post]
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

	roomId := c.Params("id")

	roomSnap, err := db.Collection("rooms").Doc(roomId).Get(context.Background())
	if err != nil {
		_ = utils.SendError(c, 404, errors.New("room not found"))
		return nil
	}

	var room map[string]interface{}
	err = roomSnap.DataTo(&room)
	if err != nil {
		_ = utils.SendError(c, 500, errors.New("unable to retrieve room information"))
		return nil
	}

	player := db.Collection("rooms").Doc(roomId).Collection("players").NewDoc()
	_, err = player.Set(context.Background(), map[string]interface{}{
		"name":      body.PlayerName,
		"timestamp": firestore.ServerTimestamp,
	})
	if err != nil {
		_ = utils.SendError(c, 400, err)
		return nil
	}

	return c.JSON(rooms.RoomJoinResponse{
		RoomId:     roomId,
		RoomName:   fmt.Sprintf("%v", room["name"]),
		PlayerId:   player.ID,
		PlayerName: body.PlayerName,
	})
}

// @Summary Get players from a room
// @Tags Rooms
// @Param id path string true "Room Id"
// @Produce json
// @Success 200 {array} players.Player
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /rooms/{id}/players [get]
func getPlayers(c *fiber.Ctx) error {
	roomId := c.Params("id")
	if len(roomId) == 0 {
		_ = utils.SendError(c, 400, errors.New("the room_id is required"))
		return nil
	}

	ctx := context.Background()

	db := new(firestore.Client)
	container.Make(&db)

	_, err := db.Collection("rooms").Doc(roomId).Get(ctx)
	if err != nil {
		_ = utils.SendError(c, 404, errors.New("room not found"))
		return nil
	}

	snaps, err := db.Collection("rooms").Doc(roomId).Collection("players").Documents(ctx).GetAll()
	if err != nil {
		_ = utils.SendError(c, 500, err)
		return nil
	}
	pls := make([]players.Player, len(snaps))
	for i, snap := range snaps {
		if err := snap.DataTo(&pls[i]); err != nil {
			_ = utils.SendError(c, 500, err)
			return nil
		}
		pls[i].Id = snap.Ref.ID
	}

	return c.JSON(pls)
}

// Registrar endpoints
func Register(app *fiber.App) {
	room := app.Group("/rooms")

	room.Post("", newRoom)
	room.Post(":id/join", joinRoom)
	room.Get(":id/players", getPlayers)
}
