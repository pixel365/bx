package types

import (
	"github.com/charmbracelet/huh"

	"github.com/pixel365/bx/internal/errors"
)

type Prompt struct {
	Value string
}

func (p *Prompt) Input(title string, validator func(string) error) error {
	if title == "" || validator == nil {
		return errors.InvalidArgumentError
	}

	value := ""
	if err := huh.NewInput().
		Title(title).
		Prompt("> ").
		Value(&value).
		Validate(validator).
		Run(); err != nil {
		return err
	}

	p.Value = value
	return nil
}

func (p *Prompt) GetValue() string {
	return p.Value
}

func NewPrompt() *Prompt {
	return &Prompt{}
}
