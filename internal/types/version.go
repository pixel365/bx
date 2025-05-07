package types

type VersionFlag string

const (
	Alpha  VersionFlag = "alpha"
	Beta   VersionFlag = "beta"
	Stable VersionFlag = "stable"
)

type Versions map[string]VersionFlag
