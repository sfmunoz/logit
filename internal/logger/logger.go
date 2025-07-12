//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/sfmunoz/logit/internal/common"
	"github.com/sfmunoz/logit/internal/handler"
)

type Logger struct {
	*slog.Logger
}

func NewLogger(h slog.Handler) *Logger {
	if h == nil {
		h = handler.NewHandler()
	}
	logInner := slog.New(h)
	return &Logger{logInner}
}

func (l *Logger) With(args ...any) *Logger {
	l.Logger = l.Logger.With(args...)
	return l
}

func (l *Logger) WithGroup(name string) *Logger {
	l.Logger = l.Logger.WithGroup(name)
	return l
}

func (l *Logger) Trace(msg string, args ...any) {
	l.Log(context.Background(), common.LevelTrace, msg, args...)
}

func (l *Logger) TraceContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, common.LevelTrace, msg, args...)
}

func (l *Logger) Notice(msg string, args ...any) {
	l.Log(context.Background(), common.LevelNotice, msg, args...)
}

func (l *Logger) NoticeContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, common.LevelNotice, msg, args...)
}

func (l *Logger) Fatal(msg string, args ...any) {
	l.Log(context.Background(), common.LevelFatal, msg, args...)
	os.Exit(1)
}

func (l *Logger) FatalContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, common.LevelFatal, msg, args...)
	os.Exit(1)
}

func (l *Logger) Panic(msg string, args ...any) {
	l.Log(context.Background(), common.LevelFatal, msg, args...)
	panic(fmt.Sprint(msg, " ", args))
}

func (l *Logger) PanicContext(ctx context.Context, msg string, args ...any) {
	l.Log(ctx, common.LevelFatal, msg, args...)
	panic(fmt.Sprint(msg, " ", args))
}
