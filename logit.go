//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package logit

import (
	"github.com/sfmunoz/logit/internal/common"
	"github.com/sfmunoz/logit/internal/logger"
)

const (
	LevelTrace  = common.LevelTrace
	LevelDebug  = common.LevelDebug
	LevelInfo   = common.LevelInfo
	LevelNotice = common.LevelNotice
	LevelWarn   = common.LevelWarn
	LevelError  = common.LevelError
	LevelFatal  = common.LevelFatal

	SymbolNone        = common.SymbolNone
	SymbolUnicodeUp   = common.SymbolUnicodeUp
	SymbolUnicodeDown = common.SymbolUnicodeDown

	TplTime    = common.TplTime
	TplUptime  = common.TplUptime
	TplLevel   = common.TplLevel
	TplSource  = common.TplSource
	TplMessage = common.TplMessage
	TplAttrs   = common.TplAttrs

	UptimeAdhoc = common.UptimeAdhoc
	UptimeStd   = common.UptimeStd
)

// make sure Level* constants are in the right order
// provided that new ones have been defined by logit
func init() {
	if LevelTrace >= LevelDebug {
		panic("LevelTrace >= LevelDebug")
	}
	if LevelDebug >= LevelInfo {
		panic("LevelDebug >= LevelInfo")
	}
	if LevelInfo >= LevelNotice {
		panic("LevelInfo >= LevelNotice")
	}
	if LevelNotice >= LevelWarn {
		panic("LevelNotice >= LevelWarn")
	}
	if LevelWarn >= LevelError {
		panic("LevelWarn >= LevelError")
	}
	if LevelError >= LevelFatal {
		panic("LevelError >= LevelFatal")
	}
}

type logit struct {
	*logger.Logger
}

func Logit() *logit {
	return &logit{logger.NewLogger(nil)}
}
