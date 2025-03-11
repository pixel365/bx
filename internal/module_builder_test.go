package internal

import (
	"fmt"
	"os"
	"testing"

	"github.com/rs/zerolog"
)

func Test_makeZipFilePath(t *testing.T) {
	mod1 := &Module{
		BuildDirectory: "testdata",
		Version:        "1.0.0",
	}

	mod2 := &Module{
		BuildDirectory: "testdata/build",
		Version:        "1.0.1",
	}

	cur, _ := os.Getwd()

	type args struct {
		module *Module
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"1", args{mod1}, fmt.Sprintf("%s/testdata/1.0.0.zip", cur), false},
		{"2", args{mod2}, fmt.Sprintf("%s/testdata/build/1.0.1.zip", cur), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := makeZipFilePath(tt.args.module)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeZipFilePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("makeZipFilePath() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeVersionDirectory(t *testing.T) {
	mod1 := &Module{
		BuildDirectory: "testdata",
		Version:        "1.0.0",
	}

	mod2 := &Module{
		BuildDirectory: "testdata/build",
		Version:        "1.0.1",
	}

	cur, _ := os.Getwd()

	type args struct {
		module *Module
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"1", args{mod1}, fmt.Sprintf("%s/testdata/1.0.0", cur), false},
		{"2", args{mod2}, fmt.Sprintf("%s/testdata/build/1.0.1", cur), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := makeVersionDirectory(tt.args.module)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeVersionDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("makeVersionDirectory() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeVersionDescription(t *testing.T) {
	mod := &Module{
		BuildDirectory: "testdata",
		Version:        "1.0.0",
	}

	type args struct {
		module *Module
		log    *zerolog.Logger
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{args{mod, nil}, "empty repository", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := makeVersionDescription(tt.args.module, tt.args.log); (err != nil) != tt.wantErr {
				t.Errorf("makeVersionDescription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
