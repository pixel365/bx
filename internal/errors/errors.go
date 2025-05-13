// Package errors provides predefined error variables used throughout the application
// to represent common failure conditions.
//
// These errors allow for consistent comparisons using errors.Is and improve clarity
// across boundaries such as services, handlers, and utilities.
package errors

import (
	"errors"
)

var (
	ErrNilModule               = errors.New("module is nil")
	ErrNoItems                 = errors.New("no items")
	ErrNoCommandSpecified      = errors.New("no command specified")
	ErrNilCookie               = errors.New("cookie is nil")
	ErrEmptySession            = errors.New("empty session")
	ErrEmptyLogin              = errors.New("empty login")
	ErrEmptyPassword           = errors.New("empty password")
	ErrPasswordTooShort        = errors.New("password is too short")
	ErrAuthentication          = errors.New("authentication failed")
	ErrEmptyModuleName         = errors.New("empty module name")
	ErrNameContainsSpace       = errors.New("name must not contain spaces")
	ErrEmptyAccountName        = errors.New("empty account name")
	ErrStageCallbackNotFound   = errors.New("stage callback not found")
	ErrChangelogValue          = errors.New("changelog from: value is required")
	ErrChangelogConditionValue = errors.New("changelog condition: value is required")
	ErrInvalidFilepath         = errors.New("invalid filepath")
	ErrSmallDepth              = errors.New("depth cannot be less than 0")
	ErrLargeDepth              = errors.New("depth cannot be greater than 5")
	ErrReplacement             = errors.New("replacement variable is empty")
	ErrEmptyVersion            = errors.New("empty version")
	ErrInvalidStages           = errors.New("stages is not valid")
	ErrInvalidRun              = errors.New("run is required")
	ErrNilRepository           = errors.New("repository is nil")
	ErrNoChanges               = errors.New("no changes detected. version directory is empty")
	ErrCallbackStage           = errors.New("callback stage is required")
	ErrCallbackPrePost         = errors.New("callback pre or post is required")
	ErrCallbackType            = errors.New("callback type is required")
	ErrCallbackMethod          = errors.New("callback method is required")
	ErrCallbackAction          = errors.New("callback action is required")
	ErrCallbackActionScheme    = errors.New(
		"callback action url scheme is invalid. allowed values are 'http' or 'https'",
	)
	ErrNilCmd                   = errors.New("cmd is nil")
	ErrNilClient                = errors.New("client is nil")
	ErrTODOContext              = errors.New("todo context is prohibited")
	ErrNilContext               = errors.New("nil context is prohibited")
	ErrInvalidRootDir           = errors.New("invalid root directory")
	ErrInvalidArgument          = errors.New("invalid argument")
	ErrDescriptionDoesNotExists = errors.New("description does not exist")
	ErrInvalidLabel             = errors.New("invalid label")
)
