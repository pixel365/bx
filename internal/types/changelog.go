package types

import (
	"golang.org/x/text/encoding/charmap"
	"strings"
)

type Changelog struct {
	From           TypeValue[ChangelogType, string]            `yaml:"from"`
	To             TypeValue[ChangelogType, string]            `yaml:"to"`
	Sort           SortingType                                 `yaml:"sort,omitempty"`
	FooterTemplate string                                      `yaml:"footerTemplate,omitempty"`
	Condition      TypeValue[ChangelogConditionType, []string] `yaml:"condition,omitempty"`
	Transform      *[]TypeValue[TransformType, []string]       `yaml:"transform,omitempty"`
}

func (c *Changelog) EncodedFooter() (string, error) {
	if c.FooterTemplate == "" {
		return "", nil
	}

	return charmap.Windows1251.NewEncoder().String("<br>" + c.FooterTemplate)
}

func (c *Changelog) ApplyTransformation(s string) string {
	if c.Transform == nil {
		return s
	}

	for _, rule := range *c.Transform {
		switch rule.Type {
		default:
			continue
		case StripPrefix:
			for _, prefix := range rule.Value {
				if strings.HasPrefix(s, prefix) {
					s = strings.TrimPrefix(s, prefix)
					break
				}
			}
		}
	}

	return strings.TrimSpace(s)
}
