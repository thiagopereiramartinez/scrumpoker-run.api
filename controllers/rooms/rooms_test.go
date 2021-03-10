package rooms

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/golobby/container"
	Assert "github.com/stretchr/testify/assert"
	"github.com/thiagopereiramartinez/scrumpoker-run.api/di"
	"golang.org/x/net/nettest"
	"io/ioutil"
	"net/http"
	"testing"
)

var app *fiber.App

func TestMain(m *testing.M) {

	_ = di.SetupDependencies()

	app = fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	Register(app)

	listener, _ := nettest.NewLocalListener("tcp")
	go func() {
		_ = app.Listener(listener)
	}()

	m.Run()

	_ = app.Shutdown()

	defer func() {
		// Encerrar conexão com o Firestore
		db := new(firestore.Client)
		container.Make(&db)
		db.Close()
	}()
}

func TestNewRoomValid(t *testing.T) {

	assert := Assert.New(t)

	body, _ := json.Marshal(map[string]interface{}{
		"name": "Room",
	})

	req, _ := http.NewRequest("POST", "/rooms", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)

	assert.NoError(err)
	assert.Equal(200, res.StatusCode)

	bodyResp, _ := ioutil.ReadAll(res.Body)
	assert.NotEmpty(string(bodyResp))
}