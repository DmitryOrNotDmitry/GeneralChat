package main

import "time"

type Message struct {
	Username  string    `json:"username" bson:"username"`
	Data      string    `json:"message" bson:"message"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}
