package logger

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

func NewLogger() zerolog.Logger {
	return zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Stack().Caller().Logger()
}

func NewLoggerWithContext(ctx context.Context) context.Context {
	zlog := NewLogger()

	return zlog.WithContext(ctx)
}
