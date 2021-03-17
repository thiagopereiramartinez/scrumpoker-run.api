package rooms

import (
	"errors"
	"strings"
	"time"
)

type Room struct {
	Id        string    `json:"id"`
	Name      string    `json:"name" firestore:"name"`
	PinCode   string    `json:"pincode" firestore:"pincode"`
	CreatedAt time.Time `json:"created_at" firestore:"timestamp"`
}

type RoomNewRequest struct {
	Name string `json:"name"`
}

func (body *RoomNewRequest) Validate() error {
	if len(strings.TrimSpace(body.Name)) == 0 {
		return errors.New("the name of the room is required")
	}

	return nil
}

type RoomNewResponse struct {
	RoomId  string `json:"room_id"`
	PinCode string `json:"pincode"`
}

type RoomJoinRequest struct {
	PlayerName string `json:"player_name"`
}

func (body *RoomJoinRequest) Validate() error {
	if len(strings.TrimSpace(body.PlayerName)) == 0 {
		return errors.New("the name of the player is required")
	}

	return nil
}

type RoomJoinResponse struct {
	Room       Room   `json:"room"`
	PlayerId   string `json:"id"`
	PlayerName string `json:"name"`
}
