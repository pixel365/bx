package internal

import (
	"github.com/rs/zerolog"
)

type Builder interface {
	Build(last bool) error
	Prepare(log *zerolog.Logger) error
	Cleanup(log *zerolog.Logger) error
	Rollback(log *zerolog.Logger) error
	Collect(last bool, log *zerolog.Logger) error
}
