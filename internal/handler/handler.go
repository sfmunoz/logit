//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package handler

import (
	"context"
	"fmt"
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

func (a attr) String() string {
	return fmt.Sprintf("%v %v %v", a.Attr, a.Value.Kind(), a.withGroup)
}

type Handler struct {
	attrs []attr

	mu       sync.Mutex
	out      io.Writer
	tsStart  time.Time
	handlers []slog.Handler

	level       slog.Leveler
	timeFormat  string
	colorMode   common.ColorMode
	colorObj    *color.Color
	symbolSet   common.SymbolSet
	tpl         []common.Tpl
	uptimeFmt   common.UptimeFormat
	attrsMode   common.AttrsMode
	replaceAttr common.ReplaceAttr
}

func NewHandler() *Handler {
	// time.RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	// time.StampMilli = "Jan _2 15:04:05.000"
	// 999: drops trailing 0; 000: keeps trailing 0
	colorMode := LogitColorModeEnv()
	return &Handler{
		attrs:      make([]attr, 0),
		out:        LogitWriterEnv(),
		tsStart:    time.Now().UTC(),
		handlers:   make([]slog.Handler, 0),
		level:      LogitLevelEnv(),
		timeFormat: LogitTimeFormatEnv(),
		colorMode:  colorMode,
		colorObj:   color.NewColor(colorMode),
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
		attrsMode:   LogitAttrsModeEnv(),
		replaceAttr: nil,
	}
}

func (h *Handler) clone() *Handler {
	ret := &Handler{
		attrs:       make([]attr, len(h.attrs)),
		out:         h.out,
		tsStart:     h.tsStart,
		handlers:    make([]slog.Handler, len(h.handlers)),
		level:       h.level,
		timeFormat:  h.timeFormat,
		colorMode:   h.colorMode,
		colorObj:    color.NewColor(h.colorMode),
		symbolSet:   h.symbolSet,
		tpl:         make([]common.Tpl, len(h.tpl)),
		uptimeFmt:   h.uptimeFmt,
		attrsMode:   h.attrsMode,
		replaceAttr: h.replaceAttr,
	}
	for i, a := range h.attrs {
		ret.attrs[i] = attr{
			Attr:      common.AttrCopy(a.Attr),
			withGroup: a.withGroup,
		}
	}
	copy(ret.handlers, h.handlers)
	copy(ret.tpl, h.tpl)
	return ret
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	for _, hh := range h.handlers {
		_ = hh.Handle(ctx, r.Clone())
	}
	buf := buffer.NewBuffer(h.timeFormat, h.colorObj, h.tsStart, h.symbolSet, h.uptimeFmt, h.attrsMode, h.replaceAttr)
	for _, tpl := range h.tpl {
		switch tpl {
		case common.TplTime:
			buf.PushTime(&r)
		case common.TplUptime:
			buf.PushUptime(&r)
		case common.TplLevel:
			buf.PushLevel(&r)
		case common.TplSource:
			buf.PushSource(&r)
		case common.TplMessage:
			buf.PushMessage(&r)
		case common.TplAttrs:
			buf.PushAttrBuiltin(&r)
			g0 := make([][]any, 0)
			g1 := []any{common.RootGroup}
			for _, a := range h.attrs {
				if a.withGroup {
					g0 = append(g0, g1)
					g1 = []any{a.Key}
				} else {
					g1 = append(g1, a.Attr)
				}
			}
			r.Attrs(func(a slog.Attr) bool {
				g1 = append(g1, a)
				return true
			})
			g0 = append(g0, g1)
			var gRoot *slog.Attr = nil
			for i := len(g0) - 1; i >= 0; i-- {
				attrs := g0[i]
				if gRoot != nil {
					attrs = append(attrs, *gRoot)
				}
				gtmp := slog.Group(attrs[0].(string), attrs[1:]...)
				gRoot = &gtmp
			}
			buf.PushAttr(gRoot)
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
	for i, a := range h.attrs {
		hc.attrs[i] = attr{Attr: common.AttrCopy(a.Attr), withGroup: a.withGroup}
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
	for i, a := range h.attrs {
		hc.attrs[i] = attr{Attr: common.AttrCopy(a.Attr), withGroup: a.withGroup}
	}
	hc.attrs[l] = attr{Attr: slog.Group(name), withGroup: true}
	return hc
}
