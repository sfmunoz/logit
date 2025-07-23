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
	sd: {lt: "ⓣ", ld: "ⓓ", li: "ⓘ ", ln: "ⓝ", lw: "ⓦ", le: "ⓔ", lf: "ⓕ"},
}

type Buffer struct {
	arr         []string
	timeFmt     string
	col         *color.Color
	tsStart     time.Time
	symbolSet   common.SymbolSet
	attrsMode   common.AttrsMode
	uptimeFmt   common.UptimeFormat
	replaceAttr common.ReplaceAttr
}

func NewBuffer(timeFmt string, col *color.Color, tsStart time.Time, symbolSet common.SymbolSet, uptimeFmt common.UptimeFormat, attrsMode common.AttrsMode, replaceAttr common.ReplaceAttr) *Buffer {
	return &Buffer{
		arr:         make([]string, 0, 20),
		timeFmt:     timeFmt,
		col:         col,
		tsStart:     tsStart,
		symbolSet:   symbolSet,
		uptimeFmt:   uptimeFmt,
		attrsMode:   attrsMode,
		replaceAttr: replaceAttr,
	}
}

func (b *Buffer) Len() int {
	return len(b.arr)
}

func (b *Buffer) WriteTo(w io.Writer) (n int64, err error) {
	tot, err := w.Write([]byte(strings.Join(b.arr, " ") + "\n"))
	return int64(tot), err
}

func (b *Buffer) pushTime(r *slog.Record, attrsReq common.AttrsMode) *Buffer {
	if r.Time.IsZero() {
		return b
	}
	if attrsReq != b.attrsMode {
		return b
	}
	switch attrsReq {
	case common.AttrsStd:
		b.arr = append(
			b.arr,
			b.col.TimFunc[0](r.Level)+r.Time.Format(b.timeFmt)+b.col.TimFunc[1](r.Level),
		)
	case common.AttrsBuiltin:
		a := slog.Time(slog.TimeKey, r.Time)
		b.PushAttr(&a)
	}
	return b
}

func (b *Buffer) PushTime(r *slog.Record) *Buffer {
	return b.pushTime(r, common.AttrsStd)
}

func (b *Buffer) PushUptime(r *slog.Record) *Buffer {
	if r.Time.IsZero() {
		return b
	}
	b.arr = append(
		b.arr,
		b.col.UptFunc[0](r.Level)+b.uptimeStr(r.Time.UTC().Sub(b.tsStart))+b.col.UptFunc[1](r.Level),
	)
	return b
}

func (b *Buffer) pushLevel(r *slog.Record, attrsReq common.AttrsMode) *Buffer {
	if attrsReq != b.attrsMode {
		return b
	}
	switch attrsReq {
	case common.AttrsStd:
		b.arr = append(
			b.arr,
			b.col.LvlFunc[0](r.Level)+lMap[b.symbolSet][r.Level]+b.col.LvlFunc[1](r.Level),
		)
	case common.AttrsBuiltin:
		a := slog.Any(slog.LevelKey, r.Level)
		b.PushAttr(&a)
	}
	return b
}

func (b *Buffer) PushLevel(r *slog.Record) *Buffer {
	return b.pushLevel(r, common.AttrsStd)
}

func (b *Buffer) pushSource(r *slog.Record, attrsReq common.AttrsMode) *Buffer {
	if attrsReq != b.attrsMode {
		return b
	}
	s := rec2src(r)
	if s == nil {
		return b
	}
	dir, file := filepath.Split(s.File)
	sourceKey := fmt.Sprintf("%s:%d", filepath.Join(filepath.Base(dir), file), s.Line)
	switch attrsReq {
	case common.AttrsStd:
		b.arr = append(
			b.arr,
			fmt.Sprintf("%s<%s>%s", b.col.SrcFunc[0](r.Level), sourceKey, b.col.SrcFunc[1](r.Level)),
		)
	case common.AttrsBuiltin:
		a := slog.Any(slog.SourceKey, sourceKey)
		b.PushAttr(&a)
	}
	return b
}

func (b *Buffer) PushSource(r *slog.Record) *Buffer {
	return b.pushSource(r, common.AttrsStd)
}

func (b *Buffer) pushMessage(r *slog.Record, attrsReq common.AttrsMode) *Buffer {
	if attrsReq != b.attrsMode {
		return b
	}
	switch attrsReq {
	case common.AttrsStd:
		b.arr = append(
			b.arr,
			b.col.MsgFunc[0](r.Level)+r.Message+b.col.MsgFunc[1](r.Level),
		)
	case common.AttrsBuiltin:
		a := slog.Any(slog.MessageKey, r.Message)
		b.PushAttr(&a)
	}
	return b
}

func (b *Buffer) PushMessage(r *slog.Record) *Buffer {
	return b.pushMessage(r, common.AttrsStd)
}

func (b *Buffer) PushAttrBuiltin(r *slog.Record) {
	b.pushTime(r, common.AttrsBuiltin)
	b.pushLevel(r, common.AttrsBuiltin)
	b.pushSource(r, common.AttrsBuiltin)
	b.pushMessage(r, common.AttrsBuiltin)
}

func (b *Buffer) PushAttr(attr *slog.Attr) *Buffer {
	attr.Value = attr.Value.Resolve()
	if attr.Equal(slog.Attr{}) || attr.Key == "" {
		return b
	}
	// https://github.com/sfmunoz/logit/issues/15 - handler.ReplaceAttr(): make use of groups
	// https://github.com/golang/example/blob/master/slog-handler-guide/README.md#implementing-handler-methods ->
	// -> https://pkg.go.dev/log/slog#HandlerOptions.ReplaceAttr
	if b.replaceAttr != nil {
		// naive support: groups are ignored -> it's OK for now despite #15 issue
		if attrNew := b.replaceAttr(nil, *attr); !attrNew.Equal(slog.Attr{}) {
			attr = &attrNew
		}
	}
	kind := attr.Value.Kind()
	if kind == slog.KindGroup {
		if attr.Key == common.RootGroup {
			for _, a := range attr.Value.Group() {
				b.PushAttr(&a)
			}
		} else {
			for _, a := range attr.Value.Group() {
				b.PushAttr(&slog.Attr{Key: attr.Key + "." + a.Key, Value: a.Value})
			}
		}
		return b
	}
	key := attr.Key
	val := attr.Value
	val2 := func() string {
		switch kind {
		case slog.KindAny:
			switch cv := val.Any().(type) {
			case encoding.TextMarshaler:
				data, err := cv.MarshalText()
				if err != nil {
					return fmt.Sprintf("!cv.MarshalText() error: %s", err)
				}
				return string(data)
			default:
				return fmt.Sprintf("%+v", val.Any())
			}
		case slog.KindBool:
			return fmt.Sprintf("%t", val.Bool())
		case slog.KindDuration:
			return b.uptimeStr(val.Duration())
		case slog.KindFloat64:
			return fmt.Sprintf("%g", val.Float64())
		case slog.KindInt64:
			return fmt.Sprintf("%d", val.Int64())
		case slog.KindString:
			return val.String()
		case slog.KindTime:
			return val.Time().String()
		case slog.KindUint64:
			return fmt.Sprintf("%d", val.Uint64())
		case slog.KindGroup:
			return "!error: slog.KindGroup already handled"
		case slog.KindLogValuer:
			return fmt.Sprintf("%+v", val.Any())
		}
		return ""
	}()
	fk := b.col.KeyFunc
	fv := b.col.ValFunc
	if _, ok := val.Any().(error); ok {
		fk = b.col.ErKFunc
		fv = b.col.ErVFunc
	}
	b.arr = append(b.arr, fk[0]()+key+"="+fk[1]()+fv[0]()+val2+fv[1]())
	return b
}

func (b *Buffer) uptimeStr(uptime time.Duration) string {
	if b.uptimeFmt == common.UptimeStd {
		return uptime.String()
	}
	timeDay := 24 * time.Hour
	days := uptime / timeDay
	uptime -= days * timeDay
	hours := uptime / time.Hour
	uptime -= hours * time.Hour
	mins := uptime / time.Minute
	uptime -= mins * time.Minute
	secs := uptime / time.Second
	uptime -= secs * time.Second
	msecs := uptime / time.Millisecond
	return fmt.Sprintf("%dd%02dh%02dm%02d.%03ds", days, hours, mins, secs, msecs)
}

func rec2src(r *slog.Record) *slog.Source {
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
