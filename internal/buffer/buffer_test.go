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
const timePat = `2[0-9]{3}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])T([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]\.[0-9]{3}Z`

var colorOff = color.NewColor(common.ColorOff)
var tsStart = time.Now().UTC()

func record(msg string, args ...any) *slog.Record {
	var pc uintptr
	var pcs [1]uintptr
	runtime.Callers(3, pcs[:])
	pc = pcs[0]
	r := slog.NewRecord(time.Now(), logit.LevelInfo, msg, pc)
	r.Add(args...)
	return &r
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

func simpleBuf(replaceAttr common.ReplaceAttr) *buffer.Buffer {
	return buffer.NewBuffer(timeFmt, colorOff, tsStart, common.SymbolNone, common.UptimeStd, common.AttrsStd, replaceAttr)
}

func TestBuffer1(t *testing.T) {
	r := record("hello")
	buf := simpleBuf(nil).PushMessage(r).PushLevel(r)
	assert(t, buf, `^hello \[I]$`)
}

func TestBuffer2(t *testing.T) {
	r := record("hello")
	assert(t, simpleBuf(nil).PushLevel(r).PushMessage(r), `^\[I] hello$`)
}

func TestBuffer3(t *testing.T) {
	r := record("hello")
	buf := simpleBuf(nil).PushSource(r).PushLevel(r).PushMessage(r)
	assert(t, buf, `^<.+/.+\.go:[0-9]+> \[I] hello$`)
}

func TestBuffer4(t *testing.T) {
	r := record("hello")
	a := slog.Int("k1", 111)
	buf := simpleBuf(nil).PushLevel(r).PushMessage(r).PushAttr(&a)
	assert(t, buf, `^\[I] hello k1=111$`)
}

func TestBuffer5(t *testing.T) {
	r := record("hello")
	buf := simpleBuf(nil).PushLevel(r).PushMessage(r).PushTime(r)
	assert(t, buf, `^\[I] hello `+timePat+`$`)
}

func TestBuffer6(t *testing.T) {
	repAttr := func(groups []string, a slog.Attr) slog.Attr {
		return slog.Attr{
			Key:   "rep-" + a.Key,
			Value: a.Value,
		}
	}
	r := record("hello")
	a1 := slog.Int("k1", 111)
	a2 := slog.Int("k2", 222)
	buf := simpleBuf(repAttr).PushLevel(r).PushMessage(r).PushTime(r).PushAttr(&a1).PushAttr(&a2)
	assert(t, buf, `^\[I] hello `+timePat+` rep-k1=111 rep-k2=222$`)
}
