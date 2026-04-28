package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

func InitLogger(withTime bool) {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
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
				source := a.Value.Any().(*slog.Source)

				dir := filepath.Base(filepath.Dir(source.File))
				file := filepath.Base(source.File)
				return slog.String(a.Key, fmt.Sprintf("%s/%s:%d", dir, file, source.Line))
			}
			return a
		},
	})
	slog.SetDefault(slog.New(handler))
}
