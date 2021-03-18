package players

import "time"

type Player struct {
	Id       string    `json:"id"`
	Name     string    `json:"name" firestore:"name"`
	Host     bool      `json:"host" firestore:"host"`
	JoinedAt time.Time `json:"joined_at" firestore:"timestamp"`
}
