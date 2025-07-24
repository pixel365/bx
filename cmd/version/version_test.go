package version

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pixel365/bx/internal/helpers"
)

func Test_newVersionCommand(t *testing.T) {
	cmd := NewVersionCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "version", cmd.Use)
	assert.True(t, cmd.HasAlias("v"))
	assert.Equal(t, "Print the version information", cmd.Short)
	assert.True(t, cmd.HasFlags())
}

func TestNewVersionCommand(t *testing.T) {
	cmd := NewVersionCommand()
	output := helpers.CaptureOutput(func() {
		_ = cmd.Execute()
	})

	want := fmt.Sprintf("%s\n", buildVersion)

	assert.Equal(t, want, output)
}

func TestNewVersionCommand_Verbose(t *testing.T) {
	cmd := NewVersionCommand()
	cmd.SetArgs([]string{"--verbose"})
	output := helpers.CaptureOutput(func() {
		_ = cmd.Execute()
	})

	want := fmt.Sprintf("Version: %s\nCommit: %s\nDate: %s\nGo: %s %s/%s\n",
		buildVersion,
		buildCommit,
		buildDate,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH)

	assert.Equal(t, want, output)
}
