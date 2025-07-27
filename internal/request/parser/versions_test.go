package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"golang.org/x/net/html"
)

func Test_versionRow(t *testing.T) {
	t.Parallel()
	v, s := versionRow(&html.Node{})
	assert.Empty(t, s)
	assert.Empty(t, v)
}

func Test_extractVersion(t *testing.T) {
	t.Parallel()
	s := extractVersion(&html.Node{})
	assert.Empty(t, s)
}

func Test_extractLabel(t *testing.T) {
	t.Parallel()
	l := extractLabel(&html.Node{})
	assert.Empty(t, l)
}
