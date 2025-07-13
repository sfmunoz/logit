//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package buffer

import (
	"encoding"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sfmunoz/logit/internal/color"
	"github.com/sfmunoz/logit/internal/common"
)

const (
	lt = common.LevelTrace
	ld = common.LevelDebug
	li = common.LevelInfo
	ln = common.LevelNotice
	lw = common.LevelWarn
	le = common.LevelError
	lf = common.LevelFatal

	sn = common.SymbolNone
	su = common.SymbolUnicodeUp
	sd = common.SymbolUnicodeDown
)

var lMap = map[common.SymbolSet]map[slog.Level]string{
	sn: {lt: "[T]", ld: "[D]", li: "[I]", ln: "[N]", lw: "[W]", le: "[E]", lf: "[F]"},
	su: {lt: "Ⓣ", ld: "Ⓓ", li: "Ⓘ", ln: "Ⓝ", lw: "Ⓦ", le: "Ⓔ", lf: "Ⓕ"},
	sd: {lt: "ⓣ", ld: "ⓓ", li: "ⓓ", ln: "ⓝ", lw: "ⓦ", le: "ⓔ", lf: "ⓕ"},
}

type Buffer struct {
	arr       []string
	timeFmt   string
	col       *color.Color
	tsStart   time.Time
	groups    []string
	symbolSet common.SymbolSet
}

func NewBuffer(timeFmt string, col *color.Color, tsStart time.Time, groups []string, symbolSet common.SymbolSet) *Buffer {
	return &Buffer{
		arr:       make([]string, 0, 20),
		timeFmt:   timeFmt,
		col:       col,
		tsStart:   tsStart,
		groups:    groups,
		symbolSet: symbolSet,
	}
}

func (b *Buffer) Len() int {
	return len(b.arr)
}

func (b *Buffer) WriteTo(w io.Writer) (n int64, err error) {
	tot, err := w.Write([]byte(strings.Join(b.arr, " ") + "\n"))
	return int64(tot), err
}

func (b *Buffer) PushTime(r slog.Record) {
	if r.Time.IsZero() {
		return
	}
	b.arr = append(
		b.arr,
		b.col.TimFunc[0](r.Level)+r.Time.Format(b.timeFmt)+b.col.TimFunc[1](r.Level),
	)
}

func (b *Buffer) PushUptime(r slog.Record) {
	if r.Time.IsZero() {
		return
	}
	b.arr = append(
		b.arr,
		b.col.UptFunc[0](r.Level)+dur2str(r.Time.UTC().Sub(b.tsStart), true)+b.col.UptFunc[1](r.Level),
	)
}

func (b *Buffer) PushLevel(r slog.Record) {
	b.arr = append(
		b.arr,
		b.col.LvlFunc[0](r.Level)+lMap[b.symbolSet][r.Level]+b.col.LvlFunc[1](r.Level),
	)
}

func (b *Buffer) PushSource(r slog.Record) {
	s := rec2src(r)
	if s == nil {
		return
	}
	dir, file := filepath.Split(s.File)
	b.arr = append(
		b.arr,
		fmt.Sprintf("%s<%s:%d>%s", b.col.SrcFunc[0](r.Level), filepath.Join(filepath.Base(dir), file), s.Line, b.col.SrcFunc[1](r.Level)),
	)
}

func (b *Buffer) PushMessage(r slog.Record) {
	b.arr = append(
		b.arr,
		b.col.MsgFunc[0](r.Level)+r.Message+b.col.MsgFunc[1](r.Level),
	)
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
	tmp := fk[0]() + pref + key + "=" + fk[1]() + fv[0]()
	switch kind {
	case slog.KindAny:
		switch cv := val.Any().(type) {
		case encoding.TextMarshaler:
			data, err := cv.MarshalText()
			if err != nil {
				tmp += fmt.Sprintf("!cv.MarshalText() error: %s", err)
				break
			}
			tmp += string(data)
		default:
			tmp += fmt.Sprintf("%+v", val.Any())
		}
	case slog.KindBool:
		tmp += fmt.Sprintf("%t", val.Bool())
	case slog.KindDuration:
		tmp += dur2str(val.Duration(), true)
	case slog.KindFloat64:
		tmp += fmt.Sprintf("%g", val.Float64())
	case slog.KindInt64:
		tmp += fmt.Sprintf("%d", val.Int64())
	case slog.KindString:
		tmp += val.String()
	case slog.KindTime:
		tmp += val.Time().String()
	case slog.KindUint64:
		tmp += fmt.Sprintf("%d", val.Uint64())
	case slog.KindGroup:
		tmp += "!error: slog.KindGroup already handled"
	case slog.KindLogValuer:
		tmp += fmt.Sprintf("%+v", val.Any())
	}
	tmp += fv[1]()
	b.arr = append(b.arr, tmp)
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
