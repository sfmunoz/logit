//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package logger_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"testing"
	"time"

	"github.com/sfmunoz/logit/internal/common"
	"github.com/sfmunoz/logit/internal/logger"
)

func TestNewLogger(t *testing.T) {
	l := logger.NewLogger(nil).
		With("a1", "v1").
		WithGroup("g1").
		WithGroup("g2").
		WithLevel(common.LevelTrace).
		WithTime(false).
		WithSource(true)
	slog.SetDefault(l.Logger)
	l.Info("logger.NewLogger()", "type", fmt.Sprintf("%T", l))
	slog.Info("Starting server", "addr", ":8080", "env", "production")
	slog.Debug("Connected to DB", "db", "myapp", "host", "localhost:5432")
	slog.Warn("Slow request", "method", "GET", "path", "/users", "duration", 497*time.Millisecond)
	slog.Error("DB connection lost", "err", "connection reset", "failure", errors.New("network off"), "db", "myapp")
	log.Print("log.Print() message")
	l.Trace("trace", "the-key", "the-val")
	l.Notice("notice", "the-key", "the-val")
	//l.Fatal("fatal", "key", "val")
	l.WithGroup("s").LogAttrs(context.Background(), common.LevelNotice, "(1) logger.WithGroup(\"s\")", slog.Int("a", 1), slog.Int("b", 2))
	l.LogAttrs(context.Background(), common.LevelNotice, "(2) logger.WithGroup(\"s\")", slog.Group("s", slog.Int("a", 1), slog.Int("b", 2)))
	slog.Log(context.Background(), common.LevelNotice, "slog.Log(LevelNotice)", "the-key", "the-val")
}
