package config

import (
	"testing"
	"time"

	"github.com/pixel365/bx/internal/model"
)

func TestNewMockConfig(t *testing.T) {
	cfg, _ := NewMockConfig()
	tests := []struct {
		want    *MockConfig
		name    string
		wantErr bool
	}{
		{cfg, "success", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMockConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMockConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want.CreatedAt.Format(time.RFC3339) != got.CreatedAt.Format(time.RFC3339) {
				t.Errorf("NewMockConfig() = %v, want %v", got, tt.want)
			}

			if tt.want.UpdatedAt.Format(time.RFC3339) != got.UpdatedAt.Format(time.RFC3339) {
				t.Errorf("NewMockConfig() = %v, want %v", got, tt.want)
			}

			if len(tt.want.Modules) != len(got.Modules) {
				t.Errorf("NewMockConfig() = %v, want %v", got, tt.want)
			}

			if len(tt.want.Accounts) != len(got.Accounts) {
				t.Errorf("NewMockConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMockConfig_GetAccounts(t *testing.T) {
	now := time.Now().UTC()
	cfg, _ := NewMockConfig()

	t.Run("success", func(t *testing.T) {
		cfg.SetAccounts(model.Account{
			CreatedAt: now,
			UpdatedAt: now,
			Login:     "test",
			Cookies:   nil,
		})

		if len(cfg.GetAccounts()) != 1 {
			t.Errorf("GetAccounts() = %v, want %v", len(cfg.GetAccounts()), 1)
		}
	})
}

func TestMockConfig_GetModules(t *testing.T) {
	now := time.Now().UTC()
	cfg, _ := NewMockConfig()

	t.Run("success", func(t *testing.T) {
		cfg.SetModules(model.Module{
			CreatedAt:   now,
			UpdatedAt:   now,
			Name:        "test",
			Path:        "test",
			Description: "test",
			Login:       "test",
		})

		if len(cfg.GetModules()) != 1 {
			t.Errorf("GetModules() = %v, want %v", len(cfg.GetModules()), 1)
		}
	})
}

func TestMockConfig_SetAccounts(t *testing.T) {
	now := time.Now().UTC()
	cfg, _ := NewMockConfig()

	t.Run("success", func(t *testing.T) {
		cfg.SetAccounts(model.Account{
			CreatedAt: now,
			UpdatedAt: now,
			Login:     "test",
			Cookies:   nil,
		})

		if len(cfg.GetAccounts()) != 1 {
			t.Errorf("GetAccounts() = %v, want %v", len(cfg.GetAccounts()), 1)
		}
	})
}

func TestMockConfig_SetModules(t *testing.T) {
	now := time.Now().UTC()
	cfg, _ := NewMockConfig()

	t.Run("success", func(t *testing.T) {
		cfg.SetModules(model.Module{
			CreatedAt:   now,
			UpdatedAt:   now,
			Name:        "test",
			Path:        "test",
			Description: "test",
			Login:       "test",
		})

		if len(cfg.GetModules()) != 1 {
			t.Errorf("GetModules() = %v, want %v", len(cfg.GetModules()), 1)
		}
	})
}

func TestMockConfig_AddAccounts(t *testing.T) {
	now := time.Now().UTC()
	cfg, _ := NewMockConfig()

	t.Run("success", func(t *testing.T) {
		cfg.AddAccounts(model.Account{
			CreatedAt: now,
			UpdatedAt: now,
			Login:     "test",
			Cookies:   nil,
		})

		if len(cfg.GetAccounts()) != 1 {
			t.Errorf("GetAccounts() = %v, want %v", len(cfg.GetAccounts()), 1)
		}

		cfg.AddAccounts(model.Account{
			CreatedAt: now,
			UpdatedAt: now,
			Login:     "test",
			Cookies:   nil,
		})

		if len(cfg.GetAccounts()) != 2 {
			t.Errorf("GetAccounts() = %v, want %v", len(cfg.GetAccounts()), 2)
		}
	})
}

func TestMockConfig_AddModules(t *testing.T) {
	now := time.Now().UTC()
	cfg, _ := NewMockConfig()

	t.Run("success", func(t *testing.T) {
		cfg.AddModules(model.Module{
			CreatedAt:   now,
			UpdatedAt:   now,
			Name:        "test",
			Path:        "test",
			Description: "test",
			Login:       "test",
		})

		if len(cfg.GetModules()) != 1 {
			t.Errorf("GetModules() = %v, want %v", len(cfg.GetModules()), 1)
		}

		cfg.AddModules(model.Module{
			CreatedAt:   now,
			UpdatedAt:   now,
			Name:        "test",
			Path:        "test",
			Description: "test",
			Login:       "test",
		})

		if len(cfg.GetModules()) != 2 {
			t.Errorf("GetModules() = %v, want %v", len(cfg.GetModules()), 2)
		}
	})
}

func TestMockConfig_Reset(t *testing.T) {
	now := time.Now().UTC()
	cfg, _ := NewMockConfig()

	t.Run("success", func(t *testing.T) {
		cfg.AddAccounts(model.Account{
			CreatedAt: now,
			UpdatedAt: now,
			Login:     "test",
			Cookies:   nil,
		})

		cfg.AddModules(model.Module{
			CreatedAt:   now,
			UpdatedAt:   now,
			Name:        "test",
			Path:        "test",
			Description: "test",
			Login:       "test",
		})

		if len(cfg.GetAccounts()) != 1 {
			t.Errorf("GetAccounts() = %v, want %v", len(cfg.GetAccounts()), 1)
		}

		if len(cfg.GetModules()) != 1 {
			t.Errorf("GetModules() = %v, want %v", len(cfg.GetModules()), 1)
		}

		err := cfg.Reset()
		if err != nil {
			t.Errorf("Reset() = %v, want nil", err)
		}

		if cfg.GetAccounts() != nil {
			t.Errorf("GetAccounts() = %v, want %v", cfg.GetAccounts(), nil)
		}

		if len(cfg.GetAccounts()) != 0 {
			t.Errorf("GetAccounts() = %v, want %v", len(cfg.GetAccounts()), 0)
		}

		if cfg.GetModules() != nil {
			t.Errorf("GetModules() = %v, want %v", cfg.GetModules(), 0)
		}

		if len(cfg.GetModules()) != 0 {
			t.Errorf("GetModules() = %v, want %v", len(cfg.GetModules()), 0)
		}
	})
}
