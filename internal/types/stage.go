package types

type Stage struct {
	Name               string           `yaml:"name"`
	To                 string           `yaml:"to"`
	ActionIfFileExists FileExistsAction `yaml:"actionIfFileExists"`
	From               []string         `yaml:"from"`
	Filter             []string         `yaml:"filter,omitempty"`
	ConvertTo1251      bool             `yaml:"convertTo1251,omitempty"`
}
