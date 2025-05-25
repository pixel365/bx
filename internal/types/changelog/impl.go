package changelog

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/charmap"

	"github.com/pixel365/bx/internal/errors"
	"github.com/pixel365/bx/internal/types"
)

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
		case types.StripPrefix:
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

func (c *Changelog) IsValid() error {
	if err := changeLogFromToValidate(c); err != nil {
		return err
	}

	if err := conditionValidate(c.Condition); err != nil {
		return err
	}

	switch c.Sort {
	case "", types.Asc, types.Desc:
	default:
		return fmt.Errorf("changelog sort must be %s or %s", types.Asc, types.Desc)
	}

	return transformValidate(c.Transform)
}

func transformValidate(transform *[]types.TypeValue[types.TransformType, []string]) error {
	if transform == nil {
		return nil
	}

	for _, rule := range *transform {
		switch rule.Type {
		default:
			return fmt.Errorf("transform rule: type must be %s", types.StripPrefix)
		case types.StripPrefix:
			for _, value := range rule.Value {
				if value == "" {
					return fmt.Errorf("transform rule: value is required")
				}
			}
		}
	}

	return nil
}

func changeLogFromToValidate(c *Changelog) error {
	if c.From.Value == "" || c.To.Value == "" {
		return errors.ErrChangelogValue
	}

	if c.From.Type != types.Commit && c.From.Type != types.Tag {
		return fmt.Errorf("changelog from: type must be %s or %s", types.Commit, types.Tag)
	}

	if c.To.Type != types.Commit && c.To.Type != types.Tag {
		return fmt.Errorf("changelog to: type must be %s or %s", types.Commit, types.Tag)
	}

	return nil
}

func conditionValidate(condition types.TypeValue[types.ChangelogConditionType, []string]) error {
	if condition.Type != "" {
		if condition.Type != types.Include &&
			condition.Type != types.Exclude {
			return fmt.Errorf(
				"changelog condition: type must be %s or %s",
				types.Include,
				types.Exclude,
			)
		}

		if len(condition.Value) == 0 {
			return errors.ErrChangelogConditionValue
		}

		for i, cond := range condition.Value {
			if cond == "" {
				return fmt.Errorf("condition [%d]: value is required", i)
			}

			_, err := regexp.Compile(cond)
			if err != nil {
				return fmt.Errorf("invalid condition [%d]: %w", i, err)
			}
		}
	}
	return nil
}
