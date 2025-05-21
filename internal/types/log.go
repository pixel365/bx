package types

type Log struct {
	Dir        string `yaml:"dir"`
	MaxSize    int    `yaml:"maxSize"`
	MaxBackups int    `yaml:"maxBackups"`
	MaxAge     int    `yaml:"maxAge"`
	LocalTime  bool   `yaml:"localTime"`
	Compress   bool   `yaml:"compress"`
}
