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
	AnsiReset  = "\033[0m"
	AnsiBold   = "\033[1m"
	AnsiFaint  = "\033[2m"
	AnsiNormal = "\033[22m"

	AnsiBlack   = "\033[30m"
	AnsiRed     = "\033[31m"
	AnsiGreen   = "\033[32m"
	AnsiYellow  = "\033[33m"
	AnsiBlue    = "\033[34m"
	AnsiMagenta = "\033[35m"
	AnsiCyan    = "\033[36m"
	AnsiWhite   = "\033[37m"
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

func cNone(...slog.Level) string {
	return ""
}

var cMap = map[slog.Level]string{
	common.LevelTrace:  AnsiCyan + AnsiFaint,
	common.LevelDebug:  AnsiWhite + AnsiFaint,
	common.LevelInfo:   AnsiGreen,
	common.LevelNotice: AnsiBlue,
	common.LevelWarn:   AnsiYellow,
	common.LevelError:  AnsiRed,
	common.LevelFatal:  AnsiRed + AnsiBold,
}

func cOn(s string) colFunc {
	if s != "" {
		return func(...slog.Level) string {
			return s
		}
	}
	return func(level ...slog.Level) string {
		if len(level) < 1 {
			return ""
		}
		return cMap[level[0]]
	}
}

func cOff(...slog.Level) string {
	return AnsiReset
}

func NewColor(cMode common.ColorMode) *Color {
	if cMode == common.ColorOff {
		return &Color{
			TimFunc: []colFunc{cNone, cNone},
			UptFunc: []colFunc{cNone, cNone},
			LvlFunc: []colFunc{cNone, cNone},
			SrcFunc: []colFunc{cNone, cNone},
			MsgFunc: []colFunc{cNone, cNone},
			KeyFunc: []colFunc{cNone, cNone},
			ValFunc: []colFunc{cNone, cNone},
			ErKFunc: []colFunc{cNone, cNone},
			ErVFunc: []colFunc{cNone, cNone},
		}
	}
	// ColorSmart
	c := &Color{
		TimFunc: []colFunc{cOn(AnsiWhite + AnsiFaint), cOff},
		UptFunc: []colFunc{cOn(AnsiCyan + AnsiFaint), cOff},
		LvlFunc: []colFunc{cOn(""), cOff},
		SrcFunc: []colFunc{cOn(AnsiWhite + AnsiFaint), cOff},
		MsgFunc: []colFunc{cOn(AnsiWhite), cOff},
		KeyFunc: []colFunc{cOn(AnsiWhite + AnsiFaint), cOff},
		ValFunc: []colFunc{cOn(AnsiWhite), cOff},
		ErKFunc: []colFunc{cOn(AnsiRed), cOff},
		ErVFunc: []colFunc{cOn(AnsiRed + AnsiBold), cOff},
	}
	// ColorMedium / ColorFull
	switch cMode {
	case common.ColorMedium:
		c.SrcFunc = []colFunc{cOn(""), cOff}
		c.MsgFunc = []colFunc{cOn(""), cOff}
	case common.ColorFull:
		c.TimFunc = []colFunc{cOn(""), cOff}
		c.UptFunc = []colFunc{cOn(""), cOff}
		c.SrcFunc = []colFunc{cOn(""), cOff}
		c.MsgFunc = []colFunc{cOn(""), cOff}
	}
	return c
}
