package players

import "time"

type Player struct {
	Id     string    `json:"id"`
	Name   string    `json:"name" firestore:"name"`
	JoinAt time.Time `json:"join_at" firestore:"timestamp"`
}