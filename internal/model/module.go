package model

import "time"

type Module struct {
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Description string    `json:"description"`
}
