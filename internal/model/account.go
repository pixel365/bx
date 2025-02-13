package model

import (
	"net/http"
	"time"
)

type Account struct {
	CreatedAt time.Time     `json:"created_at"`
	Login     string        `json:"login"`
	Cookies   []http.Cookie `json:"cookies"`
}
