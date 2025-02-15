package model

import (
	"fmt"
	"time"
)

type Module struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Description string    `json:"description"`
	Login       string    `json:"login"`
}

func (m *Module) PrintSummary(verbose bool) {
	if verbose {
		fmt.Printf("Name: %s", m.Name)
		fmt.Printf("Created At: %s\n", m.CreatedAt.Format(time.RFC822))
		fmt.Printf("Updated At: %s\n", m.UpdatedAt.Format(time.RFC822))
		fmt.Printf("Path: %s\n", m.Path)
		fmt.Printf("Description: %s\n", m.Description)
	} else {
		fmt.Printf("%s\n", m.Name)
	}
}

func (m Module) Option() string {
	return m.Name
}
