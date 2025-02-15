package model

import (
	"fmt"
	"net/http"
	"time"
)

type Account struct {
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Login     string        `json:"login"`
	Cookies   []http.Cookie `json:"cookies"`
}

func (a *Account) PrintSummary(verbose bool) {
	if verbose {
		fmt.Printf("Login: %s", a.Login)
		fmt.Printf("Created At: %s\n", a.CreatedAt.Format(time.RFC822))
		fmt.Printf("Updated At: %s\n", a.UpdatedAt.Format(time.RFC822))
		fmt.Printf("Logged in: %t\n", a.IsLoggedIn())
	} else {
		if a.IsLoggedIn() {
			fmt.Printf("%s (logged in)\n", a.Login)
		} else {
			fmt.Printf("%s\n", a.Login)
		}
	}
}

func (a *Account) IsLoggedIn() bool {
	return len(a.Cookies) > 0
}

func (a Account) Option() string {
	return a.Login
}
