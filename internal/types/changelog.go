package types

import (
	"strings"

	"golang.org/x/text/encoding/charmap"
)

type Changelog struct {
	Transform      *[]TypeValue[TransformType, []string]       `yaml:"transform,omitempty"`
	From           TypeValue[ChangelogType, string]            `yaml:"from"`
	To             TypeValue[ChangelogType, string]            `yaml:"to"`
	Sort           SortingType                                 `yaml:"sort,omitempty"`
	FooterTemplate string                                      `yaml:"footerTemplate,omitempty"`
	Condition      TypeValue[ChangelogConditionType, []string] `yaml:"condition,omitempty"`
}

// EncodedFooter returns the FooterTemplate string encoded in Windows-1251,
// prefixed with a <br> tag. If FooterTemplate is empty, it returns an empty string.
//
// Returns:
//   - The encoded footer string, or an empty string if not set
//   - An error if encoding fails
func (c *Changelog) EncodedFooter() (string, error) {
	if c.FooterTemplate == "" {
		return "", nil
	}

	return charmap.Windows1251.NewEncoder().String("<br>" + c.FooterTemplate)
}

// ApplyTransformation applies the transformation rules defined in the Transform field
// to the input string s.
//
// Currently, only the StripPrefix rule type is supported — if the string starts with
// any of the specified prefixes, the prefix is removed.
//
// After all applicable transformations are applied, the result is trimmed of
// leading and trailing whitespace using strings.TrimSpace.
//
// If no transformations are defined, the input string is returned unchanged.
//
// Returns:
//   - The transformed and trimmed string
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
