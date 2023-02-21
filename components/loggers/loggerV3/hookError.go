package loggerV3

import (
	"github.com/rs/zerolog"
)

type ErrorHook struct {
}

func (ErrorHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level == zerolog.ErrorLevel {
		// e.Str("error", "error with something")
		loggerError.Error().Caller(3).Msg(msg)
	}
}
