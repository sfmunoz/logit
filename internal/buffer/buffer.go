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
	b.WriteString(b.col.UptFunc[0](&r.Level) + dur2Str(r.Time.UTC().Sub(b.tsStart)) + b.col.UptFunc[1](&r.Level) + " ")
}

func (b *Buffer) PushLevel(l slog.Level) {
	b.WriteString(b.col.LvlFunc[0](&l) + lMap[l] + b.col.LvlFunc[1](&l) + " ")
}

func (b *Buffer) PushSource(s *slog.Source, l *slog.Level) {
	if s == nil {
		return
	}
	dir, file := filepath.Split(s.File)
	b.Printf(
		"%s<%s:%d>%s ",
		b.col.SrcFunc[0](l),
		filepath.Join(filepath.Base(dir), file),
		s.Line,
		b.col.SrcFunc[1](l),
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
	case slog.KindString:
		b.WriteString(val.String())
	case slog.KindInt64:
		b.Printf("%d", val.Int64())
	case slog.KindUint64:
		b.Printf("%d", val.Uint64())
	case slog.KindFloat64:
		b.Printf("%g", val.Float64())
	case slog.KindBool:
		b.Printf("%t", val.Bool())
	case slog.KindDuration:
		b.WriteString(val.Duration().String())
	case slog.KindTime:
		b.WriteString(val.Time().String())
	case slog.KindAny:
		switch cv := val.Any().(type) {
		case slog.Level:
			b.PushLevel(cv)
		case encoding.TextMarshaler:
			data, err := cv.MarshalText()
			if err != nil {
				b.Printf("!cv.MarshalText() error: %s", err)
				break
			}
			b.WriteString(string(data))
		case *slog.Source:
			b.PushSource(cv, nil)
		default:
			fmt.Fprintf(b, "%+v", val.Any())
		}
	}
	b.WriteString(fv[1](nil) + " ")
}

func dur2Str(dur time.Duration) string {
	// I don't like time.Duration.String() -> time.Duration.format()
	// Python example:
	//   h,m,s,ms = int((t/3600000)%3600),int((t/60000)%60),int((t/1000)%60),int(t%1000)
	//   record.relativeCreated = "{0:02d}:{1:02d}:{2:02d}.{3:03d}".format(h,m,s,ms)
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
