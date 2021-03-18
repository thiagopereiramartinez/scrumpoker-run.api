package votings

import (
	"errors"
	"strings"
	"time"
)

type Vote struct {
	PlayerId   string    `json:"player_id"`
	PlayerName string    `json:"player_name"`
	Value      float64   `json:"value"`
	VotedAt    time.Time `json:"voted_at"`
}

type RegisterVoteRequest struct {
	PlayerId string  `json:"player_id"`
	Value    float64 `json:"value"`
}

func (r *RegisterVoteRequest) Validate() error {
	if len(strings.TrimSpace(r.PlayerId)) == 0 {
		return errors.New("the player_id field is required")
	}
	return nil
}
