package rooms

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golobby/container"
	"github.com/thiagopereiramartinez/scrumpoker-run.api/internal/models/players"
	"github.com/thiagopereiramartinez/scrumpoker-run.api/internal/models/rooms"
	"github.com/thiagopereiramartinez/scrumpoker-run.api/internal/utils"
	"math/rand"
	"time"
)

var ctx context.Context

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

	seed := rand.NewSource(time.Now().UnixNano())
	rd := rand.New(seed)
	pinCode := fmt.Sprintf("%06d", rd.Intn(999999))

	doc, _, err := db.Collection("rooms").Add(ctx, map[string]interface{}{
		"name":      body.Name,
		"pincode":   pinCode,
		"timestamp": firestore.ServerTimestamp,
	})
	if err != nil {
		_ = utils.SendError(c, 500, err)
		return nil
	}

	return c.JSON(rooms.RoomNewResponse{
		RoomId:  doc.ID,
		PinCode: pinCode,
	})
}

// @Summary Join a room
// @Tags Rooms
// @Param body body rooms.RoomJoinRequest true "Join a room"
// @Param pincode path string true "Pin Code of the Room"
// @Accept json
// @Produce json
// @Success 200 {object} rooms.RoomJoinResponse
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /rooms/{pincode}/join [post]
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

	pinCode := c.Params("pincode")

	roomSnap, err := db.Collection("rooms").Where("pincode", "==", pinCode).Limit(1).Documents(ctx).Next()
	if err != nil || !roomSnap.Exists() {
		_ = utils.SendError(c, 404, errors.New("room not found"))
		return nil
	}

	roomId := roomSnap.Ref.ID

	var room rooms.Room
	err = roomSnap.DataTo(&room)
	if err != nil {
		_ = utils.SendError(c, 500, errors.New("unable to retrieve room information"))
		return nil
	}
	room.Id = roomId

	player, _, err := db.Collection("rooms").Doc(roomId).Collection("players").Add(ctx, map[string]interface{}{
		"name":      body.PlayerName,
		"timestamp": firestore.ServerTimestamp,
	})
	if err != nil {
		_ = utils.SendError(c, 500, err)
		return nil
	}

	return c.JSON(rooms.RoomJoinResponse{
		Room:       room,
		PlayerId:   player.ID,
		PlayerName: body.PlayerName,
	})
}

// @Summary Get players from a room
// @Tags Rooms
// @Param pincode path string true "Pin Code of the Room"
// @Produce json
// @Success 200 {array} players.Player
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /rooms/{pincode}/players [get]
func getPlayers(c *fiber.Ctx) error {

	pinCode := c.Params("pincode")

	db := new(firestore.Client)
	container.Make(&db)

	roomSnap, err := db.Collection("rooms").Where("pincode", "==", pinCode).Limit(1).Documents(ctx).Next()
	if err != nil || !roomSnap.Exists() {
		_ = utils.SendError(c, 404, errors.New("room not found"))
		return nil
	}

	roomId := roomSnap.Ref.ID

	snaps, err := db.Collection("rooms").Doc(roomId).Collection("players").OrderBy("timestamp", firestore.Asc).Documents(ctx).GetAll()
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
func Register(router fiber.Router) {

	ctx = context.Background()

	room := router.Group("/rooms")

	room.Post("", newRoom)
	room.Post(":pincode/join", joinRoom)
	room.Get(":pincode/players", getPlayers)
}
