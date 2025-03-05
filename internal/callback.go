package internal

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	external = "external"
	script   = "script"
)

type Runnable interface {
	PreRun() error
	PostRun() error
}

type CallbackParameters struct {
	Type       string   `yaml:"type"`
	Action     string   `yaml:"action"`
	Method     string   `yaml:"method,omitempty"`
	Parameters []string `yaml:"parameters,omitempty"`
}

type Callback struct {
	Stage string             `yaml:"stage"`
	Pre   CallbackParameters `yaml:"pre,omitempty"`
	Post  CallbackParameters `yaml:"post,omitempty"`
}

func (c Callback) PreRun() error {
	//TODO: implementation
	return nil
}

func (c Callback) PostRun() error {
	//TODO: implementation
	return nil
}

// IsValid checks if the Callback structure is valid.
// It ensures that the stage name is provided and that either pre or post parameters exist.
// Additionally, it validates the pre and post parameters if they are set.
//
// Returns:
// - error: An error if validation fails, otherwise nil.
func (c *Callback) IsValid() error {
	if c.Stage == "" {
		return errors.New("callback stage is required")
	}

	if c.Pre.Type == "" && c.Post.Type == "" {
		return errors.New("callback pre or post is required")
	}

	if c.Pre.Type != "" {
		if err := c.Pre.IsValid(); err != nil {
			return err
		}
	}

	if c.Post.Type != "" {
		if err := c.Post.IsValid(); err != nil {
			return err
		}
	}

	return nil
}

// IsValid checks if the CallbackParameters structure is valid.
// It performs type, method, action, and parameters validation.
//
// Returns:
// - error: An error if validation fails, otherwise nil.
func (c *CallbackParameters) IsValid() error {
	if err := c.validateType(); err != nil {
		return err
	}

	if err := c.validateMethod(); err != nil {
		return err
	}

	if err := c.validateAction(); err != nil {
		return err
	}

	if err := c.validateParameters(); err != nil {
		return err
	}

	return nil
}

// validateType ensures the callback type is either "script" or "external".
//
// Returns:
// - error: An error if the type is missing or invalid, otherwise nil.

func (c *CallbackParameters) validateType() error {
	if c.Type == "" {
		return errors.New("callback type is required")
	}

	if c.Type != script && c.Type != external {
		return fmt.Errorf(
			"callback type is invalid. allowed values are '%s' or '%s'",
			script,
			external,
		)
	}

	return nil
}

// validateMethod checks the HTTP method for external callbacks.
// It ensures that the method is either GET or POST when the callback type is external.
//
// Returns:
// - error: An error if the method is missing or invalid for external callbacks, otherwise nil.
func (c *CallbackParameters) validateMethod() error {
	if c.Type == external {
		if c.Method == "" {
			return errors.New("callback method is required")
		}

		if c.Method != http.MethodGet && c.Method != http.MethodPost {
			return fmt.Errorf(
				"callback method is invalid. allowed values are '%s' or '%s'",
				http.MethodGet,
				http.MethodPost,
			)
		}
	}

	return nil
}

// validateAction verifies that the action field is set and,
// if it's an external callback, checks that it's a properly formatted URL.
//
// Returns:
// - error: An error if the action is missing or improperly formatted, otherwise nil.
func (c *CallbackParameters) validateAction() error {
	if c.Action == "" {
		return errors.New("callback action is required")
	}

	if c.Type == external {
		u, err := url.Parse(c.Action)
		if err != nil {
			return fmt.Errorf("callback action url is invalid: %w", err)
		}

		if u.Scheme == "" {
			return fmt.Errorf("callback action url scheme is required")
		}

		if u.Scheme != "http" && u.Scheme != "https" {
			return errors.New(
				"callback action url scheme is invalid. allowed values are 'http' or 'https'",
			)
		}
	}

	return nil
}

// validateParameters ensures all parameters are properly formatted.
// If the callback type is external, each parameter must follow the key=value format.
//
// Returns:
// - error: An error if validation fails, otherwise nil.
func (c *CallbackParameters) validateParameters() error {
	if len(c.Parameters) == 0 {
		return nil
	}

	for _, param := range c.Parameters {
		if param == "" {
			return fmt.Errorf("callback parameter must have a value")
		}
	}

	if c.Type == external {
		for _, param := range c.Parameters {
			parts := strings.SplitN(param, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("callback parameter must have key=value format")
			}

			key, value := parts[0], parts[1]
			if key == "" {
				return fmt.Errorf("callback parameter must have a key")
			}

			if value == "" {
				return fmt.Errorf("callback parameter must have a value")
			}
		}
	}

	return nil
}
