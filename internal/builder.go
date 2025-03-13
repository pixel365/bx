package internal

import (
	"github.com/rs/zerolog"
)

type Builder interface {
	Build() error
	Prepare(log *zerolog.Logger) error
	Cleanup(log *zerolog.Logger) error
	Rollback(log *zerolog.Logger) error
	Collect(log *zerolog.Logger) error
}
