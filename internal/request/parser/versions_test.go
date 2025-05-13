package parser

import (
	"testing"

	"golang.org/x/net/html"
)

func Test_versionRow(t *testing.T) {
	t.Run("versionRow", func(t *testing.T) {
		v, s := versionRow(&html.Node{})
		if v != "" || s != "" {
			t.Errorf("got %v, %v; want empty string", v, s)
		}
	})
}

func Test_extractVersion(t *testing.T) {
	t.Run("extractVersion", func(t *testing.T) {
		s := extractVersion(&html.Node{})
		if s != "" {
			t.Errorf("got %v; want empty string", s)
		}
	})
}

func Test_extractLabel(t *testing.T) {
	t.Run("extractLabel", func(t *testing.T) {
		l := extractLabel(&html.Node{})
		if l != "" {
			t.Errorf("got %v; want empty string", l)
		}
	})
}
