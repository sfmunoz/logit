//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package logger

import (
	"io"
	"log/slog"

	"github.com/sfmunoz/logit/internal/common"
	"github.com/sfmunoz/logit/internal/handler"
)

func (l *Logger) WithWriter(w io.Writer) *Logger {
	if h, ok := l.Logger.Handler().(*handler.Handler); ok {
		return NewLogger(h.WithWriter(w))
	}
	return l
}

func (l *Logger) WithLevel(level slog.Level) *Logger {
	if h, ok := l.Logger.Handler().(*handler.Handler); ok {
		return NewLogger(h.WithLevel(level))
	}
	return l
}

func (l *Logger) WithLeveler(level slog.Leveler) *Logger {
	if h, ok := l.Logger.Handler().(*handler.Handler); ok {
		return NewLogger(h.WithLeveler(level))
	}
	return l
}

func (l *Logger) WithTimeFormat(t string) *Logger {
	if h, ok := l.Logger.Handler().(*handler.Handler); ok {
		return NewLogger(h.WithTimeFormat(t))
	}
	return l
}

func (l *Logger) WithUptimeFormat(uptimeFmt common.UptimeFormat) *Logger {
	if h, ok := l.Logger.Handler().(*handler.Handler); ok {
		return NewLogger(h.WithUptimeFormat(uptimeFmt))
	}
	return l
}

func (l *Logger) WithColorMode(cm common.ColorMode) *Logger {
	if h, ok := l.Logger.Handler().(*handler.Handler); ok {
		return NewLogger(h.WithColorMode(cm))
	}
	return l
}

func (l *Logger) WithColor(colorOn bool) *Logger {
	if h, ok := l.Logger.Handler().(*handler.Handler); ok {
		return NewLogger(h.WithColor(colorOn))
	}
	return l
}

func (l *Logger) WithHandlers(handlers []slog.Handler) *Logger {
	if h, ok := l.Logger.Handler().(*handler.Handler); ok {
		return NewLogger(h.WithHandlers(handlers))
	}
	return l
}

func (l *Logger) WithSymbolSet(symbolSet common.SymbolSet) *Logger {
	if h, ok := l.Logger.Handler().(*handler.Handler); ok {
		return NewLogger(h.WithSymbolSet(symbolSet))
	}
	return l
}

func (l *Logger) WithTpl(tpl ...common.Tpl) *Logger {
	if h, ok := l.Logger.Handler().(*handler.Handler); ok {
		return NewLogger(h.WithTpl(tpl...))
	}
	return l
}

func (l *Logger) WithReplaceAttr(replaceAttr common.ReplaceAttr) *Logger {
	if h, ok := l.Logger.Handler().(*handler.Handler); ok {
		return NewLogger(h.WithReplaceAttr(replaceAttr))
	}
	return l
}
