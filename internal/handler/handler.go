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

	addSource  bool
	level      slog.Leveler
	timeFormat string
	timeOn     bool
	uptime     bool
	colorObj   *color.Color
	symbolSet  common.SymbolSet
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
		addSource:  false,
		level:      common.LevelInfo,
		timeFormat: "2006-01-02T15:04:05.000Z07:00",
		timeOn:     true,
		uptime:     true,
		colorObj:   color.NewColor(common.ColorSmart),
		symbolSet:  common.SymbolNone,
	}
}

func (h *Handler) clone() *Handler {
	return &Handler{
		attrs:      h.attrs,  // no clone intended
		groups:     h.groups, // no clone intended
		out:        h.out,
		tsStart:    h.tsStart,
		handlers:   h.handlers, // no clone intended
		addSource:  h.addSource,
		level:      h.level,
		timeOn:     h.timeOn,
		timeFormat: h.timeFormat,
		uptime:     h.uptime,
		colorObj:   h.colorObj, // no clone intended
		symbolSet:  h.symbolSet,
	}
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	for _, hh := range h.handlers {
		_ = hh.Handle(ctx, r.Clone())
	}
	buf := buffer.NewBuffer(h.timeFormat, h.colorObj, h.tsStart, h.groups, h.symbolSet)
	defer buf.Release()
	if h.timeOn {
		buf.PushTime(r)
	}
	if h.uptime {
		buf.PushUptime(r)
	}
	buf.PushLevel(r)
	if h.addSource {
		buf.PushSource(r)
	}
	buf.PushMessage(r)
	for _, attr := range h.attrs {
		buf.PushAttr(attr)
	}
	r.Attrs(func(attr slog.Attr) bool {
		buf.PushAttr(attr)
		return true
	})
	if buf.Len() == 0 {
		return nil
	}
	buf.WriteString("\n")
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
