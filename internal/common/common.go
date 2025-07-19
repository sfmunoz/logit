//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package common

import "log/slog"

type ColorMode int
type SymbolSet int
type Tpl int
type UptimeFormat int
type ReplaceAttr func(groups []string, a slog.Attr) slog.Attr

const (
	LevelTrace  = slog.Level(-8)  // new
	LevelDebug  = slog.LevelDebug // -4
	LevelInfo   = slog.LevelInfo  // 0
	LevelNotice = slog.Level(2)   // new
	LevelWarn   = slog.LevelWarn  // 4
	LevelError  = slog.LevelError // 8
	LevelFatal  = slog.Level(12)  // new

	ColorOff ColorMode = iota
	ColorSmart
	ColorMedium
	ColorFull

	SymbolNone SymbolSet = iota
	SymbolUnicodeUp
	SymbolUnicodeDown

	TplTime Tpl = iota
	TplUptime
	TplLevel
	TplSource
	TplMessage
	TplAttrs

	UptimeAdhoc UptimeFormat = iota
	UptimeStd

	RootGroup = "__root__"
)

func AttrCopy(attr slog.Attr) slog.Attr {
	_copy := func(val slog.Value) slog.Value {
		switch val.Kind() {
		case slog.KindGroup:
			gv := val.Group()
			g := make([]slog.Attr, len(gv))
			for i, attr := range gv {
				g[i] = AttrCopy(attr)
			}
			return slog.GroupValue(g...)
		case slog.KindAny:
			any := val.Any()
			switch v := any.(type) {
			case []slog.Attr:
				g := make([]slog.Attr, len(v))
				for i, attr := range v {
					g[i] = AttrCopy(attr)
				}
				return slog.AnyValue(g)
			default:
				return val // shallow
			}
		default:
			return val
		}
	}
	return slog.Attr{
		Key:   attr.Key,
		Value: _copy(attr.Value),
	}
}
