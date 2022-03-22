package models

import "time"

type Message struct {
	HubName string    `json:"hubName"`
	Data    string    `json:"data"`
	Date    time.Time `json:"date"`
}
