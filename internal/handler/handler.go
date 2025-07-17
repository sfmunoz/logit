//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package handler

import (
	"context"
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/sfmunoz/logit/internal/buffer"
	"github.com/sfmunoz/logit/internal/color"
	"github.com/sfmunoz/logit/internal/common"
)

type attr struct {
	slog.Attr
	withGroup bool // comes from WithGroup()
}

type Handler struct {
	attrs []attr

	mu       sync.Mutex
	out      io.Writer
	tsStart  time.Time
	handlers []slog.Handler

	level       slog.Leveler
	timeFormat  string
	colorObj    *color.Color
	symbolSet   common.SymbolSet
	tpl         []common.Tpl
	uptimeFmt   common.UptimeFormat
	replaceAttr common.ReplaceAttr
}

func NewHandler() *Handler {
	// time.RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	// time.StampMilli = "Jan _2 15:04:05.000"
	// 999: drops trailing 0; 000: keeps trailing 0
	return &Handler{
		attrs:      make([]attr, 0),
		out:        LogitWriterEnv(),
		tsStart:    time.Now().UTC(),
		handlers:   make([]slog.Handler, 0),
		level:      LogitLevelEnv(),
		timeFormat: LogitTimeFormatEnv(),
		colorObj:   color.NewColor(LogitColorModeEnv()),
		symbolSet:  LogitSymbolSetEnv(),
		tpl: []common.Tpl{
			common.TplTime,
			common.TplUptime,
			common.TplLevel,
			//common.TplSource,
			common.TplMessage,
			common.TplAttrs,
		},
		uptimeFmt:   LogitUptimeFormatEnv(),
		replaceAttr: nil,
	}
}

func (h *Handler) clone() *Handler {
	return &Handler{
		attrs:       h.attrs, // no clone intended
		out:         h.out,
		tsStart:     h.tsStart,
		handlers:    h.handlers, // no clone intended
		level:       h.level,
		timeFormat:  h.timeFormat,
		colorObj:    h.colorObj, // no clone intended
		symbolSet:   h.symbolSet,
		tpl:         h.tpl, // no clone intended
		uptimeFmt:   h.uptimeFmt,
		replaceAttr: h.replaceAttr,
	}
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	for _, hh := range h.handlers {
		_ = hh.Handle(ctx, r.Clone())
	}
	buf := buffer.NewBuffer(h.timeFormat, h.colorObj, h.tsStart, h.symbolSet, h.uptimeFmt, h.replaceAttr)
	for _, tpl := range h.tpl {
		switch tpl {
		case common.TplTime:
			buf.PushTime(r)
		case common.TplUptime:
			buf.PushUptime(r)
		case common.TplLevel:
			buf.PushLevel(r)
		case common.TplSource:
			buf.PushSource(r)
		case common.TplMessage:
			buf.PushMessage(r)
		case common.TplAttrs:
			attrs_root := make([]slog.Attr, 0)
			for _, a := range h.attrs {
				attrs_root = append(attrs_root, a.Attr)
				// TODO
				// if a.withGroup { } else { }
			}
			r.Attrs(func(attr slog.Attr) bool {
				attrs_root = append(attrs_root, attr)
				return true
			})
			for _, a := range attrs_root {
				buf.PushAttr(a)
			}
		}
	}
	if buf.Len() == 0 {
		return nil
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := buf.WriteTo(h.out)
	return err
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	hc := h.clone()
	l := len(h.attrs)
	hc.attrs = make([]attr, l+len(attrs))
	if l > 0 {
		copy(hc.attrs, h.attrs)
	}
	for i, a := range attrs {
		hc.attrs[l+i] = attr{Attr: a, withGroup: false}
	}
	return hc
}

func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	hc := h.clone()
	l := len(h.attrs)
	hc.attrs = make([]attr, l+1)
	if l > 0 {
		copy(hc.attrs, h.attrs)
	}
	hc.attrs[l] = attr{Attr: slog.Group(name), withGroup: true}
	return hc
}
