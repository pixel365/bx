package client

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	t.Parallel()
	client := NewClient(10 * time.Second)
	assert.NotNil(t, client)
}
