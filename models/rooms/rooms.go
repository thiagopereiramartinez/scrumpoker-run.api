package rooms

import "errors"

type RoomNewRequest struct {
	Name string `json:"name"`
}

func (body *RoomNewRequest) Validate() error {
	if len(body.Name) == 0 {
		return errors.New("the name of the room is required")
	}

	return nil
}

type RoomNewResponse struct {
	RoomId string `json:"room_id"`
}

type RoomJoinRequest struct {
	PlayerName   string `json:"player_name"`
}

func (body *RoomJoinRequest) Validate() error {
	if len(body.PlayerName) == 0 {
		return errors.New("the name of the player is required")
	}

	return nil
}

type RoomJoinResponse struct {
	RoomId     string `json:"room_id"`
	RoomName   string `json:"room_name"`
	PlayerId   string `json:"player_id"`
	PlayerName string `json:"player_name"`
}