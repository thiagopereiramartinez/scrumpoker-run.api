package rooms

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/golobby/container"
	Assert "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thiagopereiramartinez/scrumpoker-run.api/internal/di"
	"github.com/thiagopereiramartinez/scrumpoker-run.api/internal/models"
	"github.com/thiagopereiramartinez/scrumpoker-run.api/internal/models/players"
	"github.com/thiagopereiramartinez/scrumpoker-run.api/internal/models/rooms"
	"github.com/thiagopereiramartinez/scrumpoker-run.api/test"
	"golang.org/x/net/nettest"
	"io/ioutil"
	"math/rand"
	"net/http"
	"testing"
)

var app *fiber.App
var db *firestore.Client

func TestMain(m *testing.M) {

	ctx = context.Background()

	_ = di.SetupDependencies()
	container.Make(&db)

	app = fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	Register(app)

	listener, _ := nettest.NewLocalListener("tcp")
	go func() {
		_ = app.Listener(listener)
	}()

	m.Run()

	defer func() {
		db.Close()
		_ = app.Shutdown()
	}()
}

func TestNewRoomValid(t *testing.T) {

	assert := Assert.New(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name": "Room",
	})

	req, _ := http.NewRequest("POST", "/rooms", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, 30000)

	assert.NoError(err)
	assert.Equal(200, res.StatusCode)

	bodyResp, _ := ioutil.ReadAll(res.Body)
	response := new(rooms.RoomNewResponse)
	err = json.Unmarshal(bodyResp, response)
	assert.NoError(err)

	assert.NotEmpty(response.RoomId)
	assert.NotEmpty(response.PinCode)

	snap, err := db.Collection("rooms").Doc(response.RoomId).Get(ctx)
	assert.NoError(err)
	assert.True(snap.Exists())
	assert.Equal("Room", snap.Data()["name"])
	assert.Equal(response.PinCode, snap.Data()["pincode"])
	assert.NotEmpty(snap.Data()["timestamp"])
}

func TestNewRoomNameIsEmpty(t *testing.T) {

	assert := Assert.New(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name": "",
	})

	req, _ := http.NewRequest("POST", "/rooms", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, 30000)

	assert.NoError(err)
	assert.Equal(400, res.StatusCode)

	bodyResp, _ := ioutil.ReadAll(res.Body)
	jsonBodyResp, _ := models.Error{
		Code:    400,
		Message: "the name of the room is required",
	}.ToJson()
	assert.Equal(jsonBodyResp, string(bodyResp))
}

func TestNewRoomNameIsBlank(t *testing.T) {

	assert := Assert.New(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name": "   ",
	})

	req, _ := http.NewRequest("POST", "/rooms", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, 30000)

	assert.NoError(err)
	assert.Equal(400, res.StatusCode)

	bodyResp, _ := ioutil.ReadAll(res.Body)
	jsonBodyResp, _ := models.Error{
		Code:    400,
		Message: "the name of the room is required",
	}.ToJson()
	assert.Equal(jsonBodyResp, string(bodyResp))
}

func TestNewRoomInvalidJsonBody(t *testing.T) {

	assert := Assert.New(t)

	body, _ := json.Marshal(map[string]interface{}{
		"foo": "bar",
	})

	req, _ := http.NewRequest("POST", "/rooms", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, 30000)

	assert.NoError(err)
	assert.Equal(400, res.StatusCode)

	bodyResp, _ := ioutil.ReadAll(res.Body)
	jsonBodyResp, _ := models.Error{
		Code:    400,
		Message: "the name of the room is required",
	}.ToJson()
	assert.Equal(jsonBodyResp, string(bodyResp))
}

func TestNewRoomInvalidContentType(t *testing.T) {

	assert := Assert.New(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name": "Room",
	})

	req, _ := http.NewRequest("POST", "/rooms", bytes.NewReader(body))
	req.Header.Set("Content-Type", "plain/text")
	res, err := app.Test(req, 30000)

	assert.NoError(err)
	assert.Equal(500, res.StatusCode)

	bodyResp, _ := ioutil.ReadAll(res.Body)
	jsonBodyResp, _ := models.Error{
		Code:    500,
		Message: "Unprocessable Entity",
	}.ToJson()
	assert.Equal(jsonBodyResp, string(bodyResp))
}

func TestJoinRoomValid(t *testing.T) {

	assert := Assert.New(t)

	pinCode := fmt.Sprintf("%06d", rand.Intn(999999))

	// Create a new room
	roomDoc := db.Collection("rooms").NewDoc()
	_, err := roomDoc.Set(ctx, map[string]interface{}{
		"name":      "Room",
		"pincode":   pinCode,
		"timestamp": firestore.ServerTimestamp,
	})
	assert.NoError(err)

	roomId := roomDoc.ID
	assert.NotEmpty(roomId)

	// Join a room
	body, _ := json.Marshal(rooms.RoomJoinRequest{
		PlayerName: "Thiago",
	})

	req, _ := http.NewRequest("POST", fmt.Sprintf("/rooms/%s/join", pinCode), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, 30000)

	assert.NoError(err)
	assert.Equal(200, res.StatusCode)

	bodyResp, _ := ioutil.ReadAll(res.Body)
	var result = new(rooms.RoomJoinResponse)

	err = json.Unmarshal(bodyResp, result)
	assert.NoError(err)

	roomSnap, _ := roomDoc.Get(ctx)
	var room = new(rooms.Room)
	_ = roomSnap.DataTo(room)

	assert.Equal(result.Room.Id, roomDoc.ID)
	assert.Equal(result.Room.Name, room.Name)
	assert.Equal(result.Room.PinCode, room.PinCode)
	assert.Equal(result.Room.CreatedAt, room.CreatedAt)

	playerSnap, _ := db.Collection("rooms").Doc(roomId).Collection("players").Doc(result.PlayerId).Get(ctx)
	assert.True(playerSnap.Exists())

	assert.Equal(result.PlayerName, playerSnap.Data()["name"])
}

func TestJoinRoomNameIsEmpty(t *testing.T) {

	assert := Assert.New(t)
	roomId := utils.UUID()

	body, _ := json.Marshal(rooms.RoomJoinRequest{
		PlayerName: "",
	})

	req, _ := http.NewRequest("POST", fmt.Sprintf("/rooms/%s/join", roomId), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, 30000)

	assert.NoError(err)
	assert.Equal(400, res.StatusCode)

	bodyResp, _ := ioutil.ReadAll(res.Body)
	jsonBodyResp, _ := models.Error{
		Code:    400,
		Message: "the name of the player is required",
	}.ToJson()
	assert.Equal(jsonBodyResp, string(bodyResp))
}

func TestJoinRoomNameIsBlank(t *testing.T) {

	assert := Assert.New(t)
	roomId := utils.UUID()

	body, _ := json.Marshal(rooms.RoomJoinRequest{
		PlayerName: "    ",
	})

	req, _ := http.NewRequest("POST", fmt.Sprintf("/rooms/%s/join", roomId), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, 30000)

	assert.NoError(err)
	assert.Equal(400, res.StatusCode)

	bodyResp, _ := ioutil.ReadAll(res.Body)
	jsonBodyResp, _ := models.Error{
		Code:    400,
		Message: "the name of the player is required",
	}.ToJson()
	assert.Equal(jsonBodyResp, string(bodyResp))
}

func TestJoinRoomInvalidJsonBody(t *testing.T) {

	assert := Assert.New(t)
	roomId := utils.UUID()

	body, _ := json.Marshal(map[string]interface{}{
		"foo": "bar",
	})

	req, _ := http.NewRequest("POST", fmt.Sprintf("/rooms/%s/join", roomId), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, 30000)

	assert.NoError(err)
	assert.Equal(400, res.StatusCode)

	bodyResp, _ := ioutil.ReadAll(res.Body)
	jsonBodyResp, _ := models.Error{
		Code:    400,
		Message: "the name of the player is required",
	}.ToJson()
	assert.Equal(jsonBodyResp, string(bodyResp))
}

func TestJoinRoomInvalidContentType(t *testing.T) {

	assert := Assert.New(t)
	roomId := utils.UUID()

	body, _ := json.Marshal(rooms.RoomJoinRequest{
		PlayerName: "Thiago",
	})

	req, _ := http.NewRequest("POST", fmt.Sprintf("/rooms/%s/join", roomId), bytes.NewReader(body))
	req.Header.Set("Content-Type", "plain/text")
	res, err := app.Test(req, 30000)

	assert.NoError(err)
	assert.Equal(500, res.StatusCode)

	bodyResp, _ := ioutil.ReadAll(res.Body)
	jsonBodyResp, _ := models.Error{
		Code:    500,
		Message: "Unprocessable Entity",
	}.ToJson()
	assert.Equal(jsonBodyResp, string(bodyResp))
}

func TestJoinRoomThatRoomNotExists(t *testing.T) {

	assert := Assert.New(t)
	pinCode := fmt.Sprintf("%06d", rand.Intn(999999))

	body, _ := json.Marshal(rooms.RoomJoinRequest{
		PlayerName: "Thiago",
	})

	req, _ := http.NewRequest("POST", fmt.Sprintf("/rooms/%s/join", pinCode), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req, 30000)

	assert.NoError(err)
	assert.Equal(404, res.StatusCode)

	bodyResp, _ := ioutil.ReadAll(res.Body)
	jsonBodyResp, _ := models.Error{
		Code:    404,
		Message: "room not found",
	}.ToJson()
	assert.Equal(jsonBodyResp, string(bodyResp))
}

func TestGetPlayersFromRoom(t *testing.T) {

	assert := Assert.New(t)

	// Create a new room
	pinCode := fmt.Sprintf("%06d", rand.Intn(999999))
	roomDoc := db.Collection("rooms").NewDoc()
	_, err := roomDoc.Set(ctx, map[string]interface{}{
		"name":      "Room",
		"pincode":   pinCode,
		"timestamp": firestore.ServerTimestamp,
	})
	assert.NoError(err)

	roomId := roomDoc.ID
	assert.NotEmpty(roomId)

	// Add Players
	pls := make(map[string]*players.Player, 0)
	for i := range []int{0, 1, 2, 3, 4} {
		doc := db.Collection("rooms").Doc(roomDoc.ID).Collection("players").NewDoc()
		_, err = doc.Set(ctx, map[string]interface{}{
			"name":      fmt.Sprintf("Player #%d", i),
			"timestamp": firestore.ServerTimestamp,
		})
		assert.NoError(err)

		snap, _ := db.Collection("rooms").Doc(roomDoc.ID).Collection("players").Doc(doc.ID).Get(ctx)
		pls[doc.ID] = new(players.Player)
		err = snap.DataTo(pls[doc.ID])
		pls[doc.ID].Id = doc.ID
	}

	// Get players
	req, _ := http.NewRequest("GET", fmt.Sprintf("/rooms/%s/players", pinCode), nil)
	res, err := app.Test(req, 30000)

	assert.NoError(err)
	assert.Equal(200, res.StatusCode)

	bodyResp, _ := ioutil.ReadAll(res.Body)
	respPlayers := make([]*players.Player, 0)
	_ = json.Unmarshal(bodyResp, &respPlayers)

	assert.Equal(len(pls), len(respPlayers))

	for _, p := range respPlayers {
		assert.Equal(p.Id, pls[p.Id].Id)
		assert.Equal(p.Name, pls[p.Id].Name)
		assert.Equal(p.JoinedAt, pls[p.Id].JoinedAt)
	}
}

func TestGetPlayersThatRoomNotExists(t *testing.T) {

	assert := Assert.New(t)
	pinCode := fmt.Sprintf("%06d", rand.Intn(999999))

	req, _ := http.NewRequest("GET", fmt.Sprintf("/rooms/%s/players", pinCode), nil)
	res, err := app.Test(req, 30000)

	assert.NoError(err)
	assert.Equal(404, res.StatusCode)

	bodyResp, _ := ioutil.ReadAll(res.Body)
	jsonBodyResp, _ := models.Error{
		Code:    404,
		Message: "room not found",
	}.ToJson()
	assert.Equal(jsonBodyResp, string(bodyResp))
}

func TestRegisterRoutes(t *testing.T) {

	_ = Assert.New(t)

	router := new(test.MockRouter)
	router.On("Group", "/rooms", mock.Anything).Return(router)
	router.On("Post", "", mock.Anything).Return(router)
	router.On("Post", ":pincode/join", mock.Anything).Return(router)
	router.On("Get", ":pincode/players", mock.Anything).Return(router)

	Register(router)

	router.AssertExpectations(t)

}
