package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

const (
	reset   = "\033[0m"
	cyan    = "\033[36m"
	magenta = "\033[35m"
	yellow  = "\033[33m"
	red     = "\033[31m"
	boldRed = "\033[31;1m"
)

type ColorHandler struct {
	slog.Handler
}

func (h *ColorHandler) Handle(ctx context.Context, r slog.Record) error {
	levelColor := reset

	switch r.Level {
	case slog.LevelDebug:
		levelColor = cyan
	case slog.LevelInfo:
		levelColor = magenta
	case slog.LevelWarn:
		levelColor = yellow
	case slog.LevelError:
		levelColor = red
	default:
		levelColor = reset
	}
	os.Stdout.Write([]byte(levelColor))
	defer os.Stdout.Write([]byte(reset))
	return h.Handler.Handle(ctx, r)
}

func InitLogger(withTime bool) {
	opts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				if !withTime {
					return slog.Attr{}
				}
				return slog.String(a.Key, a.Value.Time().Format("15:04:05"))
			}
			if a.Key == slog.SourceKey {
				if source, ok := a.Value.Any().(*slog.Source); ok {
					dir := filepath.Base(filepath.Dir(source.File))
					file := filepath.Base(source.File)
					return slog.String(a.Key, fmt.Sprintf("%s/%s:%d", dir, file, source.Line))
				}
			}
			return a
		}}

	baseHandler := slog.NewTextHandler(os.Stdout, opts)

	var finalHandler slog.Handler = baseHandler

	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		finalHandler = &ColorHandler{Handler: baseHandler}
	}

	slog.SetDefault(slog.New(finalHandler))
}
