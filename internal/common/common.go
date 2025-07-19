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
