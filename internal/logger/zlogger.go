package logger

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func NewLogger() zerolog.Logger {
	return zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.Kitchen}).With().Stack().Caller().Logger()
}

func NewLoggerWithContext(ctx context.Context) context.Context {
	zlog := NewLogger()

	return zlog.WithContext(ctx)
}
