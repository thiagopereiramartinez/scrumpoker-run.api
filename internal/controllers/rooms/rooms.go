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

// @Summary Get room informations
// @Tags Rooms
// @Param pincode path string true "Pin Code"
// @Produce json
// @Success 200 {object} rooms.RoomGetResponse
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /rooms/{pincode} [get]
func getRoom(c *fiber.Ctx) error {

	pinCode := c.Params("pincode")

	db := new(firestore.Client)
	container.Make(&db)

	roomSnaps, err := db.Collection("rooms").Where("pincode", "==", pinCode).Limit(1).Documents(ctx).GetAll()
	if err != nil || len(roomSnaps) == 0 {
		_ = utils.SendError(c, 404, errors.New("room not found"))
		return nil
	}
	roomSnap := roomSnaps[0]

	var room rooms.RoomGetResponse
	room.Id = roomSnap.Ref.ID
	room.Name = roomSnap.Data()["name"].(string)
	room.PinCode = roomSnap.Data()["pincode"].(string)

	if topic := roomSnap.Data()["topic"]; topic != nil {
		room.Topic = topic.(string)
	}
	room.CreatedAt = roomSnap.Data()["timestamp"].(time.Time)

	snaps, err := db.Collection("rooms").Doc(room.Id).Collection("players").OrderBy("timestamp", firestore.Asc).Documents(ctx).GetAll()
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

	room.Players = pls

	return c.JSON(room)
}

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
// @Param pincode path string true "Pin Code"
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

	roomSnaps, err := db.Collection("rooms").Where("pincode", "==", pinCode).Limit(1).Documents(ctx).GetAll()
	if err != nil || len(roomSnaps) == 0 {
		_ = utils.SendError(c, 404, errors.New("room not found"))
		return nil
	}
	roomSnap := roomSnaps[0]

	roomId := roomSnap.Ref.ID

	var room rooms.Room
	err = roomSnap.DataTo(&room)
	if err != nil {
		_ = utils.SendError(c, 500, errors.New("unable to retrieve room information"))
		return nil
	}
	room.Id = roomId

	var playerHost = false
	if body.PlayerHost != nil && *body.PlayerHost == true {
		if snaps, err := db.Collection("rooms").Doc(roomId).Collection("players").Where("host", "==", true).Documents(ctx).GetAll(); err != nil || len(snaps) > 0 {
			_ = utils.SendError(c, 400, errors.New("there is already a user defined as host"))
			return nil
		}
		playerHost = true
	}

	player, _, err := db.Collection("rooms").Doc(roomId).Collection("players").Add(ctx, map[string]interface{}{
		"name":      body.PlayerName,
		"host":      playerHost,
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
		PlayerHost: playerHost,
	})
}

// @Summary Set a topic
// @Tags Rooms
// @Param body body rooms.RoomTopicRequest true "Set a topic"
// @Param pincode path string true "Pin Code"
// @Accept json
// @Success 200 {string} string "OK"
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /rooms/{pincode}/topic [put]
func setTopic(c *fiber.Ctx) error {

	body := new(rooms.RoomTopicRequest)
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

	roomSnaps, err := db.Collection("rooms").Where("pincode", "==", pinCode).Limit(1).Documents(ctx).GetAll()
	if err != nil || len(roomSnaps) == 0 {
		_ = utils.SendError(c, 404, errors.New("room not found"))
		return nil
	}
	roomSnap := roomSnaps[0]

	_, err = db.Collection("rooms").Doc(roomSnap.Ref.ID).Set(ctx, map[string]interface{}{
		"topic": body.TopicName,
	}, firestore.MergeAll)
	if err != nil {
		_ = utils.SendError(c, 500, errors.New("unable to update room information"))
		return nil
	}

	return c.SendStatus(200)
}

// @Summary Get players from a room
// @Tags Rooms
// @Param pincode path string true "Pin Code"
// @Produce json
// @Success 200 {array} players.Player
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /rooms/{pincode}/players [get]
func getPlayers(c *fiber.Ctx) error {

	pinCode := c.Params("pincode")

	db := new(firestore.Client)
	container.Make(&db)

	roomSnaps, err := db.Collection("rooms").Where("pincode", "==", pinCode).Limit(1).Documents(ctx).GetAll()
	if err != nil || len(roomSnaps) == 0 {
		_ = utils.SendError(c, 404, errors.New("room not found"))
		return nil
	}
	roomSnap := roomSnaps[0]

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

// @Summary Rename player
// @Tags Rooms
// @Param body body rooms.RenameUserRequest true "User name"
// @Param pincode path string true "Pin Code"
// @Param id path string true "Player ID"
// @Accept json
// @Success 200 {string} string "OK"
// @Failure 400 {object} models.Error
// @Failure 404 {object} models.Error
// @Failure 500 {object} models.Error
// @Router /rooms/{pincode}/players/{id} [patch]
func renameUser(c *fiber.Ctx) error {

	pinCode := c.Params("pincode")
	playerId := c.Params("id")

	body := new(rooms.RenameUserRequest)
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

	roomSnaps, err := db.Collection("rooms").Where("pincode", "==", pinCode).Limit(1).Documents(ctx).GetAll()
	if err != nil || len(roomSnaps) == 0 {
		_ = utils.SendError(c, 404, errors.New("room not found"))
		return nil
	}
	roomSnap := roomSnaps[0]

	roomId := roomSnap.Ref.ID

	playerSnap, err := db.Collection("rooms").Doc(roomId).Collection("players").Doc(playerId).Get(ctx)
	if err != nil || !playerSnap.Exists() {
		_ = utils.SendError(c, 404, errors.New("player not found"))
		return nil
	}

	_, _ = db.Collection("rooms").Doc(roomId).Collection("players").Doc(playerId).Update(ctx, []firestore.Update {
		{
			Path: "name",
			Value: body.PlayerName,
		},
	})

	return c.SendStatus(200)
}

// Register endpoints
func Register(router fiber.Router) {

	ctx = context.Background()

	room := router.Group("/rooms")

	room.Post("", newRoom)
	room.Get(":pincode", getRoom)
	room.Post(":pincode/join", joinRoom)
	room.Put(":pincode/topic", setTopic)
	room.Get(":pincode/players", getPlayers)
	room.Patch(":pincode/players/:id", renameUser)
}
