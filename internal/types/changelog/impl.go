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
// Supported rule types:
//   - StripPrefix: removes the prefix if the string starts with any of the specified values.
//   - StripSuffix: removes the suffix if the string ends with any of the specified values.
//
// After applying all applicable transformations, the result is trimmed of
// leading and trailing whitespace using strings.TrimSpace.
//
// If no transformations are defined, the input string is returned unchanged.
//
// Returns the transformed and trimmed string.
func (c *Changelog) ApplyTransformation(s string) string {
	if c.Transform == nil {
		return s
	}

	for _, rule := range *c.Transform {
		switch rule.Type {
		default:
			continue
		case types.StripPrefix:
			s = stripPrefix(s, rule.Value)
		case types.StripSuffix:
			s = stripSuffix(s, rule.Value)
		case types.RemoveAll:
			s = removeAll(s, rule.Value)
		}
	}

	s = truncate(s, c.MaxLength)

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

	if c.MaxLength < 0 {
		return fmt.Errorf("changelog max length must be non-negative")
	}

	return transformValidate(c.Transform)
}

func transformValidate(transform *[]types.TypeValue[types.TransformType, []string]) error {
	if transform == nil {
		return nil
	}

	for _, rule := range *transform {
		if len(rule.Value) == 0 {
			return fmt.Errorf("transform rule: value is empty")
		}

		switch rule.Type {
		default:
			return fmt.Errorf("transform rule: type must be %s", types.StripPrefix)
		case types.StripPrefix, types.StripSuffix, types.RemoveAll:
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

func stripPrefix(s string, values []string) string {
	for _, value := range values {
		if strings.HasPrefix(s, value) {
			s = strings.TrimPrefix(s, value)
			break
		}
	}

	return s
}

func stripSuffix(s string, values []string) string {
	for _, value := range values {
		if strings.HasSuffix(s, value) {
			s = strings.TrimSuffix(s, value)
			break
		}
	}

	return s
}

func removeAll(s string, values []string) string {
	changed := false
	for _, value := range values {
		if strings.Contains(s, value) {
			changed = true
			s = strings.ReplaceAll(s, value, "")
		}
	}

	if changed {
		s = strings.Join(strings.Fields(s), " ")
	}

	return s
}

func truncate(s string, max int) string {
	if max <= 0 || len(s) <= max {
		return s
	}

	return s[:max]
}
