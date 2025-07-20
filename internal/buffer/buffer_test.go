//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package buffer_test

import (
	"bytes"
	"log/slog"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/sfmunoz/logit"
	"github.com/sfmunoz/logit/internal/buffer"
	"github.com/sfmunoz/logit/internal/color"
	"github.com/sfmunoz/logit/internal/common"
)

const timeFmt = "2006-01-02T15:04:05.000Z07:00"

var colorOff = color.NewColor(common.ColorOff)
var tsStart = time.Now().UTC()

func record(msg string, args ...any) slog.Record {
	var pc uintptr
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	pc = pcs[0]
	r := slog.NewRecord(time.Now(), logit.LevelInfo, msg, pc)
	r.Add(args...)
	return r
}

func assert(t *testing.T, buf *buffer.Buffer, re string) {
	want, err := regexp.Compile(re)
	if err != nil {
		t.Fatalf("regexp.Compile(%s) failed: %s", re, err)
	}
	var out bytes.Buffer
	_, err = buf.WriteTo(&out)
	if err != nil {
		t.Fatalf("buf.WriteTo() failed: %s", err)
	}
	got := out.String()
	if !strings.HasSuffix(got, "\n") {
		t.Fatalf("assert(): got='%s' doesn't have '\\n' suffix", got)
	}
	got = strings.TrimRight(got, "\n")
	if !want.MatchString(got) {
		t.Fatalf("assert(): got='%s' doesn't match want='%s'", got, want)
	}
}

func simpleBuf() *buffer.Buffer {
	return buffer.NewBuffer(timeFmt, colorOff, tsStart, common.SymbolNone, common.UptimeStd, nil)
}

func TestBuffer1(t *testing.T) {
	buf := simpleBuf()
	r := record("hello")
	buf.PushMessage(r)
	buf.PushLevel(r)
	assert(t, buf, `^hello \[I]$`)
}

func TestBuffer2(t *testing.T) {
	buf := simpleBuf()
	r := record("hello")
	buf.PushLevel(r)
	buf.PushMessage(r)
	assert(t, buf, `^\[I] hello$`)
}

func TestBuffer3(t *testing.T) {
	buf := simpleBuf()
	r := record("hello")
	buf.PushSource(r)
	buf.PushLevel(r)
	buf.PushMessage(r)
	assert(t, buf, `^<.+/.+:[0-9]+> \[I] hello$`)
}
