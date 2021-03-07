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
