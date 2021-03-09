package players

import "time"

type Player struct {
	Id       string    `json:"id"`
	Name     string    `json:"name" firestore:"name"`
	JoinedAt time.Time `json:"joined_at" firestore:"timestamp"`
}
