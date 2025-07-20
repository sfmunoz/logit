//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//
//

package logger_test

import (
	"bytes"
	"context"
	"log/slog"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/sfmunoz/logit"
	"github.com/sfmunoz/logit/internal/logger"
)

// 2025-07-20T16:41:50.744Z
const timePat = `2[0-9]{3}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])T([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]\.[0-9]{3}Z`

// 0d00h00m00.000s
const durPat = `[0-9]d([0-1][0-9]|2[0-3])h[0-5][0-9]m[0-5][0-9]\.[0-9]{3}s`

var ctx = context.Background()

func record(msg string, args ...any) slog.Record {
	var pc uintptr
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	pc = pcs[0]
	r := slog.NewRecord(time.Now(), logit.LevelInfo, msg, pc)
	r.Add(args...)
	return r
}

func assert(t *testing.T, l *logger.Logger, msg string, re string) {
	want, err := regexp.Compile(re)
	if err != nil {
		t.Fatalf("regexp.Compile(%s) failed: %s", re, err)
	}
	var out bytes.Buffer
	l.WithWriter(&out).Info(msg)
	got := out.String()
	if !strings.HasSuffix(got, "\n") {
		t.Fatalf("assert(): got='%s' doesn't have '\\n' suffix", got)
	}
	got = strings.TrimRight(got, "\n")
	if !want.MatchString(got) {
		t.Fatalf("assert(): got='%s' doesn't match want='%s'", got, want)
	}
}

func TestLogger1(t *testing.T) {
	l := logger.NewLogger(nil).
		WithColor(false)
	re := `^` + timePat + ` ` + durPat + ` \[I] hello$`
	assert(t, l, "hello", re)
}

func TestLogger2(t *testing.T) {
	l := logger.NewLogger(nil).
		WithSymbolSet(logit.SymbolUnicodeUp).
		WithColor(false)
	re := `^` + timePat + ` ` + durPat + ` Ⓘ hello$`
	assert(t, l, "hello", re)
}

func TestLogger3(t *testing.T) {
	l := logger.NewLogger(nil).
		WithColor(false).
		With("k1", "v1")
	re := `^` + timePat + ` ` + durPat + ` \[I] hello k1=v1$`
	assert(t, l, "hello", re)
}

func TestLogger4(t *testing.T) {
	l := logger.NewLogger(nil).
		WithColor(false).
		With("k1", "v1").
		WithGroup("g1").
		With("k2", "v2")
	re := `^` + timePat + ` ` + durPat + ` \[I] hello k1=v1 g1.k2=v2$`
	assert(t, l, "hello", re)
}
