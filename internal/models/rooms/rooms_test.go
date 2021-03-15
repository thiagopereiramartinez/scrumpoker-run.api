package rooms

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoomNewRequestValid(t *testing.T) {

	room := RoomNewRequest{
		Name: "thiago",
	}
	assert.NoError(t, room.Validate())

}

func TestRoomNewRequestInvalid(t *testing.T) {

	room := RoomNewRequest{
		Name: "",
	}
	assert.Errorf(t, room.Validate(), "the name of the room is required")

}

func TestRoomJoinRequestValid(t *testing.T) {

	room := RoomJoinRequest{
		PlayerName: "thiago",
	}
	assert.NoError(t, room.Validate())

}

func TestRoomJoinRequestInvalid(t *testing.T) {

	room := RoomJoinRequest{
		PlayerName: "",
	}
	assert.Errorf(t, room.Validate(), "the name of the room is required")

}
