package changelog

import (
	"github.com/pixel365/bx/internal/types"
)

type Changelog struct {
	Transform      *[]types.TypeValue[types.TransformType, []string]       `yaml:"transform,omitempty"`
	From           types.TypeValue[types.ChangelogType, string]            `yaml:"from"`
	To             types.TypeValue[types.ChangelogType, string]            `yaml:"to"`
	Sort           types.SortingType                                       `yaml:"sort,omitempty"`
	FooterTemplate string                                                  `yaml:"footerTemplate,omitempty"`
	Condition      types.TypeValue[types.ChangelogConditionType, []string] `yaml:"condition,omitempty"`
	MaxLength      int                                                     `yaml:"maxLength,omitempty"`
}
