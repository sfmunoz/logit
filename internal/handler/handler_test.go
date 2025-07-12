//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//
//

package handler_test

import (
	"errors"
	"log"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/sfmunoz/logit/internal/handler"
)

func TestRaw(t *testing.T) {
	h := handler.NewHandler()
	logger := slog.New(h)
	slog.SetDefault(logger)
	slog.Info("Starting server", "addr", ":8080", "env", "production")
	slog.Debug("Connected to DB", "db", "myapp", "host", "localhost:5432")
	slog.Warn("Slow request", "method", "GET", "path", "/users", "duration", 497*time.Millisecond)
	slog.Error("DB connection lost", "err", "connection reset", "failure", errors.New("network off"), "db", "myapp")
	log.Print("log.Print() message")
}

func TestOpts(t *testing.T) {
	h := handler.NewHandler().
		WithWriter(os.Stderr).
		WithSource(true).
		WithLevel(slog.LevelDebug)
	logger := slog.New(h)
	slog.SetDefault(logger)
	slog.Info("Starting server", "addr", ":8080", "env", "production")
	slog.Debug("Connected to DB", "db", "myapp", "host", "localhost:5432")
	slog.Warn("Slow request", "method", "GET", "path", "/users", "duration", 497*time.Millisecond)
	slog.Error("DB connection lost", "err", "connection reset", "failure", errors.New("network off"), "db", "myapp")
}

func TestFanout(t *testing.T) {
	lv := new(slog.LevelVar)
	lv.Set(slog.LevelDebug)
	opts1 := &slog.HandlerOptions{AddSource: true, Level: lv}
	h1 := slog.NewTextHandler(os.Stderr, opts1)
	opts2 := &slog.HandlerOptions{AddSource: false, Level: lv}
	h2 := slog.NewJSONHandler(os.Stderr, opts2)
	h := handler.NewHandler().
		WithLeveler(lv).
		WithHandlers([]slog.Handler{h1, h2}).
		WithTime(false)
	logger := slog.New(h)
	slog.SetDefault(logger)
	slog.Info("Message repeated", "times", 3)
}

func TestWithAttrsGroup(t *testing.T) {
	handler := handler.NewHandler()
	logger := slog.New(handler).
		With("\"a1\"", "v1").
		WithGroup("g1").
		WithGroup("g2").
		With("\"a2\"", "v2")
	slog.SetDefault(logger)
	slog.Info("Starting server", "addr", ":8080", "env", "production")
	slog.Debug("Connected to DB", "db", "myapp", "host", "localhost:5432")
	slog.Warn("Slow request", "method", "GET", "path", "/users", "duration", 497*time.Millisecond)
	slog.Error("DB connection lost", "err", "connection reset", "failure", errors.New("network off"), "db", "myapp")
	log.Print("log.Print() message")
}
