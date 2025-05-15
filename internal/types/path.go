package types

type Path struct {
	From           string
	To             string
	ActionIfExists FileExistsAction
	Convert        bool
}
