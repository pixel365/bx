package config

import (
	"sync"
	"time"

	"github.com/pixel365/bx/internal/model"
)

type Config struct {
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Accounts  []model.Account `json:"accounts,omitempty"`
	Modules   []model.Module  `json:"modules,omitempty"`
	mu        sync.RWMutex
}
