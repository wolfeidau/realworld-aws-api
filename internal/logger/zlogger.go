package logger

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewLogger() zerolog.Logger {
	return log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Kitchen}).With().Stack().Caller().Logger()
}

func NewLoggerWithContext(ctx context.Context) context.Context {
	zlog := NewLogger()

	return zlog.WithContext(ctx)
}
