package callback

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pixel365/bx/internal/errors"
	"github.com/pixel365/bx/internal/helpers"
	"github.com/pixel365/bx/internal/validators"

	"github.com/go-cmd/cmd"
)

const (
	ExternalType = "external"
	CommandType  = "command"
)

// CallbackParameters defines the details of a callback action.
//
// Fields:
//   - Type:     The type of callback (e.g., "http").
//   - Action:   The action to be triggered (e.g., a URL or command name).
//   - Method:   Optional HTTP method or execution method (e.g., "GET", "POST").
//   - Parameters: Optional list of arguments or parameters to pass during callback execution.
type CallbackParameters struct {
	Type       string   `yaml:"type"`
	Action     string   `yaml:"action"`
	Method     string   `yaml:"method,omitempty"`
	Parameters []string `yaml:"parameters,omitempty"`
}

// Callback represents a stage-specific callback definition.
//
// A callback can define separate `pre` and `post` actions to be executed before or after
// a specific stage.
//
// Fields:
//   - Stage: The name of the stage where the callback applies.
//   - Pre:   Parameters for the action to be executed before the stage.
//   - Post:  Parameters for the action to be executed after the stage.
type Callback struct {
	Stage string             `yaml:"stage"`
	Pre   CallbackParameters `yaml:"pre,omitempty"`
	Post  CallbackParameters `yaml:"post,omitempty"`
}

func (c Callback) PreRun(ctx context.Context) error {
	if err := c.Pre.IsValid(); err != nil {
		return err
	}

	if err := c.Pre.Run(ctx); err != nil {
		return fmt.Errorf("pre run callback failed for stage %s: %w", c.Stage, err)
	}

	return nil
}

func (c Callback) PostRun(ctx context.Context) error {
	if err := c.Post.IsValid(); err != nil {
		return err
	}

	if err := c.Post.Run(ctx); err != nil {
		return fmt.Errorf("post run callback failed for stage %s: %w", c.Stage, err)
	}

	return nil
}

// IsValid checks if the Callback structure is valid.
// It ensures that the stage name is provided and that either pre or post parameters exist.
// Additionally, it validates the pre and post parameters if they are set.
//
// Returns:
//   - error: An error if validation fails, otherwise nil.
func (c *Callback) IsValid() error {
	if c.Stage == "" {
		return errors.ErrCallbackStage
	}

	if c.Pre.Type == "" && c.Post.Type == "" {
		return errors.ErrCallbackPrePost
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
//   - error: An error if validation fails, otherwise nil.
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

// Run executes the callback based on its type (external or command).
// It first validates the callback parameters, then either runs an external HTTP request
// or a system command depending on the `Type` of the callback.
//
// Parameters:
//   - ctx (context.Context): The context for the execution, used for cancellation and timeouts.
//
// Returns:
//   - error: Returns an error if validation fails or if execution of the callback fails.
func (c *CallbackParameters) Run(ctx context.Context) error {
	if err := c.IsValid(); err != nil {
		return err
	}

	if c.Type == ExternalType {
		return c.runExternal(ctx)
	}

	if c.Type == CommandType {
		return c.runCommand(ctx)
	}

	return nil
}

// validateType ensures the callback type is either "CommandType" or "ExternalType".
//
// Returns:
// 	- error: An error if the type is missing or invalid, otherwise nil.

func (c *CallbackParameters) validateType() error {
	if c.Type == "" {
		return errors.ErrCallbackType
	}

	if c.Type != CommandType && c.Type != ExternalType {
		return fmt.Errorf(
			"callback type is invalid. allowed values are '%s' or '%s'",
			CommandType,
			ExternalType,
		)
	}

	return nil
}

// validateMethod checks the HTTP method for ExternalType callbacks.
// It ensures that the method is either GET or POST when the callback type is ExternalType.
//
// Returns:
//   - error: An error if the method is missing or invalid for ExternalType callbacks, otherwise nil.
func (c *CallbackParameters) validateMethod() error {
	if c.Type == ExternalType {
		if c.Method == "" {
			return errors.ErrCallbackMethod
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
// if it's an ExternalType callback, checks that it's a properly formatted URL.
//
// Returns:
//   - error: An error if the action is missing or improperly formatted, otherwise nil.
func (c *CallbackParameters) validateAction() error {
	if c.Action == "" {
		return errors.ErrCallbackAction
	}

	if c.Type == ExternalType {
		u, err := url.Parse(c.Action)
		if err != nil {
			return fmt.Errorf("callback action url is invalid: %w", err)
		}

		if u.Scheme == "" {
			return fmt.Errorf("callback action url scheme is required")
		}

		if u.Scheme != "http" && u.Scheme != "https" {
			return errors.ErrCallbackActionScheme
		}
	}

	return nil
}

// validateParameters ensures all parameters are properly formatted.
// If the callback type is ExternalType, each parameter must follow the key=value format.
//
// Returns:
//   - error: An error if validation fails, otherwise nil.
func (c *CallbackParameters) validateParameters() error {
	if len(c.Parameters) == 0 {
		return nil
	}

	for i, param := range c.Parameters {
		if param == "" {
			return fmt.Errorf("callback parameter[%d] is empty", i)
		}

		if c.Type == CommandType {
			if !validators.ValidateArgument(param) {
				return fmt.Errorf("callback parameter[%d] is invalid", i)
			}
		}
	}

	if c.Type == ExternalType {
		for i, param := range c.Parameters {
			parts := strings.SplitN(param, "=", 2)
			if len(parts) != 2 {
				return fmt.Errorf("callback parameter[%d] must have key=value format", i)
			}

			key, value := parts[0], parts[1]
			if key == "" {
				return fmt.Errorf("callback parameter[%d] must have a key", i)
			}

			if value == "" {
				return fmt.Errorf("callback parameter[%d] must have a value", i)
			}
		}
	}

	return nil
}

// runExternal executes an external HTTP request based on the callback parameters.
// It constructs the URL and body for the request, sends the request to the specified endpoint,
// and checks for a successful HTTP response status code (200 OK).
// The function has a timeout of 30 seconds for the request.
//
// Parameters:
//   - ctx (context.Context): The context for the execution, used for cancellation and timeouts.
//
// Returns:
//   - error: Returns an error if the request fails or if the response status code is not OK.
func (c *CallbackParameters) runExternal(ctx context.Context) error {
	ttlCtx, cancelFunc := context.WithTimeout(ctx, 30*time.Second)
	defer cancelFunc()

	u, body := c.buildUrlAndBody()
	req, err := http.NewRequestWithContext(ttlCtx, c.Method, u, body)
	if err != nil {
		return err
	}

	client := &http.Client{}

	//nolint:bodyclose
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer helpers.Cleanup(resp.Body, nil)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("callback returned status code %d", resp.StatusCode)
	}

	return nil
}

// runCommand executes a command defined in the callback parameters.
// It validates the command arguments, starts the command execution asynchronously,
// and logs the output (stdout and stderr). The function has a 30-second deadline for execution.
// If the command fails or the context is cancelled, it handles the errors accordingly.
//
// Parameters:
//   - ctx (context.Context): The context for execution, used for cancellation and deadlines.
//
// Returns:
//   - error: Returns an error if the command execution fails, arguments are invalid, or the context is cancelled.
func (c *CallbackParameters) runCommand(ctx context.Context) error {
	ctx, cancelFunc := context.WithDeadline(ctx, time.Now().Add(30*time.Second))
	defer cancelFunc()

	rawCommand, _ := c.buildUrlAndBody()
	args := strings.Fields(rawCommand)

	for _, arg := range args[1:] {
		if !validators.ValidateArgument(arg) {
			return fmt.Errorf("invalid callback CommandType argument '%s'", arg)
		}
	}

	com := cmd.NewCmd(args[0], args[1:]...)
	statusChan := com.Start()

	select {
	case status := <-statusChan:
		if status.Error != nil {
			return fmt.Errorf("callback CommandType failed: %w", status.Error)
		}
	case <-ctx.Done():
		if err := com.Stop(); err != nil {
			return fmt.Errorf("callback CommandType failed: %w", err)
		}
	}

	return nil
}

// buildUrlAndBody constructs the URL and the body for the callback based on its type and parameters.
// If there are parameters, they are appended to the action URL for `GET` requests or included in the body
// for other HTTP methods.
// For `CommandType`, the action and parameters are combined into a command string.
//
// Returns:
//   - string: The constructed URL or command string.
//   - io.Reader: The body of the request if applicable (nil for GET requests or CommandType).
func (c *CallbackParameters) buildUrlAndBody() (string, io.Reader) {
	if len(c.Parameters) == 0 {
		return c.Action, nil
	}

	if c.Type == ExternalType {
		if c.Method == http.MethodGet {
			delimiter := "?"
			query := strings.Join(c.Parameters, "&")
			if strings.Contains(c.Action, "?") {
				delimiter = "&"
			}

			return fmt.Sprintf("%s%s%s", c.Action, delimiter, query), nil
		}

		body := url.Values{}
		for _, param := range c.Parameters {
			parts := strings.SplitN(param, "=", 2)
			if len(parts) != 2 {
				continue
			}
			body.Add(parts[0], parts[1])
		}

		return c.Action, strings.NewReader(body.Encode())
	}

	return fmt.Sprintf("%s %s", c.Action, strings.Join(c.Parameters, " ")), nil
}

// ValidateCallbacks validates a list of Callback objects by invoking their IsValid method.
//
// Iterates through the slice of callbacks and calls IsValid on each one. If any callback
// returns a validation error, the function wraps it with the callback index for context
// and returns immediately. If all callbacks are valid or the slice is empty, returns nil.
//
// Parameters:
//   - callbacks: A slice of Callback instances to validate.
//
// Returns:
//   - error: An error if any callback fails validation, wrapped with its index; otherwise nil.
func ValidateCallbacks(callbacks []Callback) error {
	if len(callbacks) > 0 {
		for i := range callbacks {
			cb := callbacks[i]
			if err := cb.IsValid(); err != nil {
				return fmt.Errorf("callback [%d]: %w", i, err)
			}
		}
	}

	return nil
}
