package types

type FileExistsAction string
type ChangelogType string
type ChangelogConditionType string
type SortingType string
type BuildType string
type TransformType string

const (
	Replace        FileExistsAction = "replace"
	Skip           FileExistsAction = "skip"
	ReplaceIfNewer FileExistsAction = "replace_if_newer"

	Commit ChangelogType = "commit"
	Tag    ChangelogType = "tag"

	Include ChangelogConditionType = "include"
	Exclude ChangelogConditionType = "exclude"

	Asc  SortingType = "asc"
	Desc SortingType = "desc"

	StripPrefix TransformType = "stripPrefix"
	StripSuffix TransformType = "stripSuffix"
	RemoveAll   TransformType = "removeAll"
)

type TypeValue[T1 any, T2 any] struct {
	Type  T1 `yaml:"type"`
	Value T2 `yaml:"value"`
}

type Builds struct {
	Release     []string `yaml:"release"`
	LastVersion []string `yaml:"lastVersion,omitempty"`
}
