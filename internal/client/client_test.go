package client

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	t.Run("new client", func(t *testing.T) {
		client := NewClient(10 * time.Second)
		if client == nil {
			t.Error("nil client")
		}
	})
}
