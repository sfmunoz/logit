//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//
//

package buffer_test

import (
	"bytes"
	"log/slog"
	"runtime"
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

func assert(t *testing.T, buf *buffer.Buffer, want string) {
	var out bytes.Buffer
	_, err := buf.WriteTo(&out)
	if err != nil {
		t.Fatalf("buf.WriteTo() failed: %s", err)
	}
	got := out.String()
	if got != want {
		t.Fatalf("assert(): want='%s', got='%s'", want, got)
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
	assert(t, buf, "hello [I]\n")
}

func TestBuffer2(t *testing.T) {
	buf := simpleBuf()
	r := record("hello")
	buf.PushLevel(r)
	buf.PushMessage(r)
	assert(t, buf, "[I] hello\n")
}
