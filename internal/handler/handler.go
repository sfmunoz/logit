//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package handler

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/sfmunoz/logit/internal/buffer"
	"github.com/sfmunoz/logit/internal/color"
	"github.com/sfmunoz/logit/internal/common"
)

type Handler struct {
	attrs  []slog.Attr
	groups []string

	mu       sync.Mutex
	out      io.Writer
	tsStart  time.Time
	handlers []slog.Handler

	level      slog.Leveler
	timeFormat string
	colorObj   *color.Color
	symbolSet  common.SymbolSet
	tpl        []common.Tpl
	uptimeFmt  common.UptimeFormat
}

func NewHandler() *Handler {
	// time.RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	// time.StampMilli = "Jan _2 15:04:05.000"
	// 999: drops trailing 0; 000: keeps trailing 0
	return &Handler{
		attrs:      make([]slog.Attr, 0),
		groups:     make([]string, 0),
		out:        os.Stderr,
		tsStart:    time.Now().UTC(),
		handlers:   make([]slog.Handler, 0),
		level:      common.LevelInfo,
		timeFormat: "2006-01-02T15:04:05.000Z07:00",
		colorObj:   color.NewColor(common.ColorSmart),
		symbolSet:  common.SymbolNone,
		tpl: []common.Tpl{
			common.TplTime,
			common.TplUptime,
			common.TplLevel,
			//common.TplSource,
			common.TplMessage,
			common.TplAttrs,
		},
		uptimeFmt: common.UptimeAdhoc,
	}
}

func (h *Handler) clone() *Handler {
	return &Handler{
		attrs:      h.attrs,  // no clone intended
		groups:     h.groups, // no clone intended
		out:        h.out,
		tsStart:    h.tsStart,
		handlers:   h.handlers, // no clone intended
		level:      h.level,
		timeFormat: h.timeFormat,
		colorObj:   h.colorObj, // no clone intended
		symbolSet:  h.symbolSet,
		tpl:        h.tpl, // no clone intended
		uptimeFmt:  h.uptimeFmt,
	}
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	for _, hh := range h.handlers {
		_ = hh.Handle(ctx, r.Clone())
	}
	buf := buffer.NewBuffer(h.timeFormat, h.colorObj, h.tsStart, h.groups, h.symbolSet, h.uptimeFmt)
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
			for _, attr := range h.attrs {
				buf.PushAttr(attr)
			}
			r.Attrs(func(attr slog.Attr) bool {
				buf.PushAttr(attr)
				return true
			})
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
	hc.attrs = attrs
	return hc
}

func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	hc := h.clone()
	hc.groups = append(hc.groups, name)
	return hc
}
