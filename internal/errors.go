package internal

import (
	"errors"
)

var (
	NilModuleError                = errors.New("module is nil")
	NoItemsError                  = errors.New("no items")
	NoCommandSpecifiedError       = errors.New("no command specified")
	NilCookieError                = errors.New("cookie is nil")
	EmptySessionError             = errors.New("empty session")
	EmptyLoginError               = errors.New("empty login")
	EmptyPasswordError            = errors.New("empty password")
	PasswordTooShortError         = errors.New("password is too short")
	AuthenticationError           = errors.New("authentication failed")
	EmptyModuleNameError          = errors.New("empty module name")
	NameContainsSpaceError        = errors.New("name must not contain spaces")
	EmptyAccountNameError         = errors.New("empty account name")
	StageCallbackNotFoundError    = errors.New("stage callback not found")
	InvalidChangelogSettingsError = errors.New("invalid changelog settings")
	ChangelogFromValueError       = errors.New("changelog from: value is required")
	ChangelogToValueError         = errors.New("changelog to: value is required")
	ChangelogConditionValueError  = errors.New("changelog condition: value is required")
	InvalidFilepathError          = errors.New("invalid filepath")
	SmallDepthError               = errors.New("depth cannot be less than 0")
	LargeDepthError               = errors.New("depth cannot be greater than 5")
	ReplacementError              = errors.New("replacement variable is empty")
	EmptyVersionError             = errors.New("empty version")
	InvalidStagesError            = errors.New("stages is not valid")
	InvalidRunError               = errors.New("run is required")
	NilRepositoryError            = errors.New("repository is nil")
	NoChangesError                = errors.New("no changes detected. version directory is empty")
	CallbackStageError            = errors.New("callback stage is required")
	CallbackPrePostError          = errors.New("callback pre or post is required")
	CallbackTypeError             = errors.New("callback type is required")
	CallbackMethodError           = errors.New("callback method is required")
	CallbackActionError           = errors.New("callback action is required")
	CallbackActionSchemeError     = errors.New(
		"callback action url scheme is invalid. allowed values are 'http' or 'https'",
	)
)
