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
	"runtime"
	"time"

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

func (l *Logger) clone() *Logger {
	c := *l
	return &c
}

func (l *Logger) With(args ...any) *Logger {
	if len(args) == 0 {
		return l
	}
	c := l.clone()
	c.Logger = l.Logger.With(args...)
	return c
}

func (l *Logger) WithGroup(name string) *Logger {
	if name == "" {
		return l
	}
	c := l.clone()
	c.Logger = l.Logger.WithGroup(name)
	return c
}

func (l *Logger) Trace(msg string, args ...any) {
	l.log(context.Background(), common.LevelTrace, msg, args...)
}

func (l *Logger) TraceContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, common.LevelTrace, msg, args...)
}

func (l *Logger) Notice(msg string, args ...any) {
	l.log(context.Background(), common.LevelNotice, msg, args...)
}

func (l *Logger) NoticeContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, common.LevelNotice, msg, args...)
}

func (l *Logger) Fatal(msg string, args ...any) {
	l.log(context.Background(), common.LevelFatal, msg, args...)
	os.Exit(1)
}

func (l *Logger) FatalContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, common.LevelFatal, msg, args...)
	os.Exit(1)
}

func (l *Logger) Panic(msg string, args ...any) {
	l.log(context.Background(), common.LevelFatal, msg, args...)
	panic(fmt.Sprint(msg, " ", args))
}

func (l *Logger) PanicContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, common.LevelFatal, msg, args...)
	panic(fmt.Sprint(msg, " ", args))
}

// log() method copied and adapted from <go1.24.5.linux-amd64>/src/log/slog/logger.go
// which has the following license
//
// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
func (l *Logger) log(ctx context.Context, level slog.Level, msg string, args ...any) {
	if !l.Enabled(ctx, level) {
		return
	}
	var pc uintptr
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	pc = pcs[0]
	r := slog.NewRecord(time.Now(), level, msg, pc)
	r.Add(args...)
	if ctx == nil {
		ctx = context.Background()
	}
	_ = l.Handler().Handle(ctx, r)
}
