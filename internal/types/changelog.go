package types

type Changelog struct {
	From      TypeValue[ChangelogType, string]            `yaml:"from"`
	To        TypeValue[ChangelogType, string]            `yaml:"to"`
	Sort      ChangelogSort                               `yaml:"sort,omitempty"`
	Condition TypeValue[ChangelogConditionType, []string] `yaml:"condition,omitempty"`
}
