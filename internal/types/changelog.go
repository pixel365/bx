package types

import "golang.org/x/text/encoding/charmap"

type Changelog struct {
	From           TypeValue[ChangelogType, string]            `yaml:"from"`
	To             TypeValue[ChangelogType, string]            `yaml:"to"`
	Sort           SortingType                                 `yaml:"sort,omitempty"`
	FooterTemplate string                                      `yaml:"footerTemplate,omitempty"`
	Condition      TypeValue[ChangelogConditionType, []string] `yaml:"condition,omitempty"`
}

func (c *Changelog) EncodedFooter() (string, error) {
	if c.FooterTemplate == "" {
		return "", nil
	}

	return charmap.Windows1251.NewEncoder().String("<br>" + c.FooterTemplate)
}
