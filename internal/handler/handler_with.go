//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package handler

import (
	"io"
	"log/slog"

	"github.com/sfmunoz/logit/internal/color"
	"github.com/sfmunoz/logit/internal/common"
)

func (h *Handler) WithWriter(out io.Writer) *Handler {
	c := h.clone()
	c.out = out
	return c
}

func (h *Handler) WithLevel(level slog.Level) *Handler {
	return h.WithLeveler(level)
}

func (h *Handler) WithLeveler(level slog.Leveler) *Handler {
	c := h.clone()
	c.level = level
	return c
}

func (h *Handler) WithTimeFormat(t string) *Handler {
	c := h.clone()
	c.timeFormat = t
	return c
}

func (h *Handler) WithUptimeFormat(uptimeFmt common.UptimeFormat) slog.Handler {
	hc := h.clone()
	hc.uptimeFmt = uptimeFmt
	return hc
}

func (h *Handler) WithColorMode(cm common.ColorMode) *Handler {
	c := h.clone()
	c.colorObj = color.NewColor(cm)
	return c
}

func (h *Handler) WithColor(colorOn bool) *Handler {
	if colorOn {
		return h.WithColorMode(common.ColorSmart)
	}
	return h.WithColorMode(common.ColorOff)
}

func (h *Handler) WithHandlers(handlers []slog.Handler) *Handler {
	c := h.clone()
	c.handlers = handlers
	return c
}

func (h *Handler) WithSymbolSet(symbolSet common.SymbolSet) *Handler {
	c := h.clone()
	c.symbolSet = symbolSet
	return c
}

func (h *Handler) WithTpl(tpl ...common.Tpl) slog.Handler {
	hc := h.clone()
	hc.tpl = tpl
	return hc
}
