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

type ColFunc func(...slog.Level) string

type Color struct {
	TimFunc []ColFunc
	UptFunc []ColFunc
	LvlFunc []ColFunc
	SrcFunc []ColFunc
	MsgFunc []ColFunc
	KeyFunc []ColFunc
	ValFunc []ColFunc
	ErKFunc []ColFunc
	ErVFunc []ColFunc
}

func cDynamic() ColFunc {
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

func cSet(s string) ColFunc {
	return func(...slog.Level) string {
		return s
	}
}

func cReset() ColFunc {
	return cSet(ansiReset)
}

func cNone() ColFunc {
	return cSet("")
}

func NewColor(cMode common.ColorMode) *Color {
	// ColorOff
	if cMode == common.ColorOff {
		return &Color{
			TimFunc: []ColFunc{cNone(), cNone()},
			UptFunc: []ColFunc{cNone(), cNone()},
			LvlFunc: []ColFunc{cNone(), cNone()},
			SrcFunc: []ColFunc{cNone(), cNone()},
			MsgFunc: []ColFunc{cNone(), cNone()},
			KeyFunc: []ColFunc{cNone(), cNone()},
			ValFunc: []ColFunc{cNone(), cNone()},
			ErKFunc: []ColFunc{cNone(), cNone()},
			ErVFunc: []ColFunc{cNone(), cNone()},
		}
	}
	// ColorSmart
	c := &Color{
		TimFunc: []ColFunc{cSet(ansiWhite + ansiFaint), cReset()},
		UptFunc: []ColFunc{cSet(ansiCyan + ansiFaint), cReset()},
		LvlFunc: []ColFunc{cDynamic(), cReset()},
		SrcFunc: []ColFunc{cSet(ansiWhite + ansiFaint), cReset()},
		MsgFunc: []ColFunc{cSet(ansiWhite), cReset()},
		KeyFunc: []ColFunc{cSet(ansiWhite + ansiFaint), cReset()},
		ValFunc: []ColFunc{cSet(ansiWhite), cReset()},
		ErKFunc: []ColFunc{cSet(ansiRed), cReset()},
		ErVFunc: []ColFunc{cSet(ansiRed + ansiBold), cReset()},
	}
	// ColorMedium / ColorFull
	switch cMode {
	case common.ColorMedium:
		c.SrcFunc = []ColFunc{cDynamic(), cReset()}
		c.MsgFunc = []ColFunc{cDynamic(), cReset()}
	case common.ColorFull:
		c.TimFunc = []ColFunc{cDynamic(), cReset()}
		c.UptFunc = []ColFunc{cDynamic(), cReset()}
		c.SrcFunc = []ColFunc{cDynamic(), cReset()}
		c.MsgFunc = []ColFunc{cDynamic(), cReset()}
	}
	return c
}
