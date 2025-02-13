package model

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/color"
)

type Account struct {
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Login     string        `json:"login"`
	Cookies   []http.Cookie `json:"cookies"`
}

func (a *Account) PrintSummary(verbose bool) {
	if verbose {
		color.Green("Login: %s", a.Login)
		fmt.Printf("Created At: %s\n", a.CreatedAt)
		fmt.Printf("Updated At: %s\n", a.UpdatedAt)
		fmt.Printf("Is Authenticated: %t\n", len(a.Cookies) > 0)
	} else {
		fmt.Printf("%s\n", a.Login)
	}
}
