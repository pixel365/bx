package types

type VersionLabel string

const (
	Alpha  VersionLabel = "alpha"
	Beta   VersionLabel = "beta"
	Stable VersionLabel = "stable"
)

type Versions map[string]VersionLabel
