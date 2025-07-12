//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package buffer

import (
	"bytes"
	"encoding"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sfmunoz/logit/internal/color"
	"github.com/sfmunoz/logit/internal/common"
)

type Buffer struct {
	*bytes.Buffer
	timeFmt string
	col     *color.Color
	tsStart time.Time
	groups  []string
}

var lMap = map[slog.Level]string{
	common.LevelTrace:  "[T]",
	common.LevelDebug:  "[D]",
	common.LevelInfo:   "[I]",
	common.LevelNotice: "[N]",
	common.LevelWarn:   "[W]",
	common.LevelError:  "[E]",
	common.LevelFatal:  "[F]",
}

var pool = sync.Pool{
	New: func() any {
		return &Buffer{
			Buffer: bytes.NewBuffer(make([]byte, 0, 1024)),
		}
	},
}

func NewBuffer(timeFmt string, col *color.Color, tsStart time.Time, groups []string) *Buffer {
	b := pool.Get().(*Buffer)
	b.timeFmt = timeFmt
	b.col = col
	b.tsStart = tsStart
	b.groups = groups
	return b
}

func (b *Buffer) Release() {
	if b.Cap() > 16384 {
		return
	}
	b.Reset()
	pool.Put(b)
}

func (b *Buffer) Printf(format string, a ...any) {
	fmt.Fprintf(b, format, a...)
}

func (b *Buffer) PushTime(r slog.Record) {
	if r.Time.IsZero() {
		return
	}
	b.WriteString(b.col.TimFunc[0](&r.Level) + r.Time.Format(b.timeFmt) + b.col.TimFunc[1](&r.Level) + " ")
}

func (b *Buffer) PushUptime(r slog.Record) {
	if r.Time.IsZero() {
		return
	}
	b.WriteString(b.col.UptFunc[0](&r.Level) + dur2str(r.Time.UTC().Sub(b.tsStart), true) + b.col.UptFunc[1](&r.Level) + " ")
}

func (b *Buffer) PushLevel(r slog.Record) {
	b.WriteString(b.col.LvlFunc[0](&r.Level) + lMap[r.Level] + b.col.LvlFunc[1](&r.Level) + " ")
}

func (b *Buffer) PushSource(r slog.Record) {
	s := rec2src(r)
	if s == nil {
		return
	}
	dir, file := filepath.Split(s.File)
	b.Printf(
		"%s<%s:%d>%s ",
		b.col.SrcFunc[0](&r.Level),
		filepath.Join(filepath.Base(dir), file),
		s.Line,
		b.col.SrcFunc[1](&r.Level),
	)
}

func (b *Buffer) PushMessage(r slog.Record) {
	b.WriteString(b.col.MsgFunc[0](&r.Level) + r.Message + b.col.MsgFunc[1](&r.Level) + " ")
}

func (b *Buffer) PushAttr(attr slog.Attr) {
	attr.Value = attr.Value.Resolve()
	if attr.Equal(slog.Attr{}) || attr.Key == "" {
		return
	}
	kind := attr.Value.Kind()
	if kind == slog.KindGroup {
		for _, a := range attr.Value.Group() {
			k := attr.Key + "." + a.Key
			b.PushAttr(slog.Attr{Key: k, Value: a.Value})
		}
		return
	}
	key := attr.Key
	val := attr.Value
	fk := b.col.KeyFunc
	fv := b.col.ValFunc
	if _, ok := val.Any().(error); ok {
		fk = b.col.ErKFunc
		fv = b.col.ErVFunc
	}
	pref := ""
	if len(b.groups) > 0 {
		pref = strings.Join(b.groups, ".") + "."
	}
	b.WriteString(fk[0](nil) + pref + key + "=" + fk[1](nil) + fv[0](nil))
	switch kind {
	case slog.KindAny:
		switch cv := val.Any().(type) {
		case encoding.TextMarshaler:
			data, err := cv.MarshalText()
			if err != nil {
				b.Printf("!cv.MarshalText() error: %s", err)
				break
			}
			b.WriteString(string(data))
		default:
			fmt.Fprintf(b, "%+v", val.Any())
		}
	case slog.KindBool:
		b.Printf("%t", val.Bool())
	case slog.KindDuration:
		b.WriteString(dur2str(val.Duration(), true))
	case slog.KindFloat64:
		b.Printf("%g", val.Float64())
	case slog.KindInt64:
		b.Printf("%d", val.Int64())
	case slog.KindString:
		b.WriteString(val.String())
	case slog.KindTime:
		b.WriteString(val.Time().String())
	case slog.KindUint64:
		b.Printf("%d", val.Uint64())
	case slog.KindGroup:
		b.Printf("!error: slog.KindGroup already handled")
	case slog.KindLogValuer:
		fmt.Fprintf(b, "%+v", val.Any())
	}
	b.WriteString(fv[1](nil) + " ")
}

func dur2str(dur time.Duration, adhoc bool) string {
	if !adhoc {
		return dur.String()
	}
	timeDay := 24 * time.Hour
	days := dur / timeDay
	dur -= days * timeDay
	hours := dur / time.Hour
	dur -= hours * time.Hour
	mins := dur / time.Minute
	dur -= mins * time.Minute
	secs := dur / time.Second
	dur -= secs * time.Second
	msecs := dur / time.Millisecond
	return fmt.Sprintf("%dd%02dh%02dm%02d.%03ds", days, hours, mins, secs, msecs)
}

func rec2src(r slog.Record) *slog.Source {
	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()
	if f.File == "" {
		return nil
	}
	return &slog.Source{
		Function: f.Function,
		File:     f.File,
		Line:     f.Line,
	}
}
