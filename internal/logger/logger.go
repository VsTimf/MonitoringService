package logger

import (
	"io"
	"log/slog"
	"time"

	"github.com/lmittmann/tint"
)

func New(w io.Writer, logLevel slog.Leveler) *slog.Logger {
	handler := tint.NewTextHandler(w, &tint.Options{
		Level:      logLevel,
		TimeFormat: time.TimeOnly,
		NoColor:    false,
	})

	return slog.New(handler)
}
