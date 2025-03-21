package cmd

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal"
)

func Test_newPushCommand(t *testing.T) {
	cmd := newPushCommand()

	t.Run("should create new push command", func(t *testing.T) {
		if cmd == nil {
			t.Error("cmd is nil")
		}

		if cmd.Use != "push" {
			t.Errorf("new push command should use 'push', got '%s'", cmd.Use)
		}

		if cmd.Short != "Push module to a Marketplace" {
			t.Errorf("cmd.Short should be 'Push module to a Marketplace' but got '%s'", cmd.Short)
		}

		if cmd.RunE == nil {
			t.Error("cmd RunE should not be nil")
		}

		if !cmd.HasFlags() {
			t.Errorf("cmd.HasFlags() should be true")
		}

		if cmd.HasSubCommands() {
			t.Errorf("cmd.HasSubCommands() should be false")
		}

		if len(cmd.Aliases) > 0 {
			t.Errorf("len(cmd.Aliases) should be 0 but got %d", len(cmd.Aliases))
		}
	})
}

func Test_handlePassword(t *testing.T) {
	cmd := newPushCommand()
	cmd.SetArgs([]string{"--password", "123456", "--name", "test"})
	_ = cmd.Flags().Set("password", "123456")
	_ = cmd.Flags().Set("name", "test")

	module := &internal.Module{}

	type args struct {
		cmd    *cobra.Command
		module *internal.Module
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"success", args{cmd, module}, "123456", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handlePassword(tt.args.cmd, tt.args.module)
			if (err != nil) != tt.wantErr {
				t.Errorf("handlePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handlePassword() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_push_nil(t *testing.T) {
	t.Run("nil command", func(t *testing.T) {
		err := push(nil, []string{})
		if err == nil {
			t.Errorf("err is nil")
		}

		if !errors.Is(err, internal.NilCmdError) {
			t.Errorf("err = %v, want %v", err, internal.NilCmdError)
		}
	})
}
