//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package color

import (
	"log/slog"

	"github.com/sfmunoz/logit/internal/common"
)

const (
	ansiReset  = "\033[0m"
	ansiBold   = "\033[1m"
	ansiFaint  = "\033[2m"
	ansiNormal = "\033[22m"

	ansiBlack   = "\033[30m"
	ansiRed     = "\033[31m"
	ansiGreen   = "\033[32m"
	ansiYellow  = "\033[33m"
	ansiBlue    = "\033[34m"
	ansiMagenta = "\033[35m"
	ansiCyan    = "\033[36m"
	ansiWhite   = "\033[37m"
)

type colFunc func(...slog.Level) string

type Color struct {
	TimFunc []colFunc
	UptFunc []colFunc
	LvlFunc []colFunc
	SrcFunc []colFunc
	MsgFunc []colFunc
	KeyFunc []colFunc
	ValFunc []colFunc
	ErKFunc []colFunc
	ErVFunc []colFunc
}

func cDynamic() colFunc {
	var m = map[slog.Level]string{
		common.LevelTrace:  ansiCyan + ansiFaint,
		common.LevelDebug:  ansiWhite + ansiFaint,
		common.LevelInfo:   ansiGreen,
		common.LevelNotice: ansiBlue,
		common.LevelWarn:   ansiYellow,
		common.LevelError:  ansiRed,
		common.LevelFatal:  ansiRed + ansiBold,
	}
	return func(level ...slog.Level) string {
		if len(level) < 1 {
			return ""
		}
		return m[level[0]]
	}
}

func cSet(s string) colFunc {
	return func(...slog.Level) string {
		return s
	}
}

func cReset() colFunc {
	return cSet(ansiReset)
}

func cNone() colFunc {
	return cSet("")
}

func NewColor(cMode common.ColorMode) *Color {
	// ColorOff
	if cMode == common.ColorOff {
		return &Color{
			TimFunc: []colFunc{cNone(), cNone()},
			UptFunc: []colFunc{cNone(), cNone()},
			LvlFunc: []colFunc{cNone(), cNone()},
			SrcFunc: []colFunc{cNone(), cNone()},
			MsgFunc: []colFunc{cNone(), cNone()},
			KeyFunc: []colFunc{cNone(), cNone()},
			ValFunc: []colFunc{cNone(), cNone()},
			ErKFunc: []colFunc{cNone(), cNone()},
			ErVFunc: []colFunc{cNone(), cNone()},
		}
	}
	// ColorSmart
	c := &Color{
		TimFunc: []colFunc{cSet(ansiWhite + ansiFaint), cReset()},
		UptFunc: []colFunc{cSet(ansiCyan + ansiFaint), cReset()},
		LvlFunc: []colFunc{cDynamic(), cReset()},
		SrcFunc: []colFunc{cSet(ansiWhite + ansiFaint), cReset()},
		MsgFunc: []colFunc{cSet(ansiWhite), cReset()},
		KeyFunc: []colFunc{cSet(ansiWhite + ansiFaint), cReset()},
		ValFunc: []colFunc{cSet(ansiWhite), cReset()},
		ErKFunc: []colFunc{cSet(ansiRed), cReset()},
		ErVFunc: []colFunc{cSet(ansiRed + ansiBold), cReset()},
	}
	// ColorMedium / ColorFull
	switch cMode {
	case common.ColorMedium:
		c.SrcFunc = []colFunc{cDynamic(), cReset()}
		c.MsgFunc = []colFunc{cDynamic(), cReset()}
	case common.ColorFull:
		c.TimFunc = []colFunc{cDynamic(), cReset()}
		c.UptFunc = []colFunc{cDynamic(), cReset()}
		c.SrcFunc = []colFunc{cDynamic(), cReset()}
		c.MsgFunc = []colFunc{cDynamic(), cReset()}
	}
	return c
}
