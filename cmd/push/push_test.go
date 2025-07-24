package push

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newPushCommand(t *testing.T) {
	cmd := NewPushCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "push", cmd.Use)
	assert.Equal(t, "Push module to a Marketplace", cmd.Short)
	assert.NotNil(t, cmd.RunE)
	assert.True(t, cmd.HasFlags())
	assert.False(t, cmd.HasSubCommands())
	assert.Empty(t, cmd.Aliases)
}
