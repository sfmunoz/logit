//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//
//

package color_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/sfmunoz/logit/internal/color"
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

func prune(s string) string {
	// naive implementation... but it's enough
	return strings.ReplaceAll(s, "\033", "\\033")
}

func colAssert(t *testing.T, msg, got, want string) {
	if got == want {
		return
	}
	t.Fatalf("%s failed: want='%+v', got='%+v'", msg, prune(want), prune(got))
}

func dynAssert(t *testing.T, name string, fn []color.ColFunc) {
	colAssert(t, fmt.Sprintf("%s[0](common.LevelTrace)", name), fn[0](common.LevelTrace), ansiCyan+ansiFaint)
	colAssert(t, fmt.Sprintf("%s[0](common.LevelDebug)", name), fn[0](common.LevelDebug), ansiWhite+ansiFaint)
	colAssert(t, fmt.Sprintf("%s[0](common.LevelInfo)", name), fn[0](common.LevelInfo), ansiGreen)
	colAssert(t, fmt.Sprintf("%s[0](common.LevelNotice)", name), fn[0](common.LevelNotice), ansiBlue)
	colAssert(t, fmt.Sprintf("%s[0](common.LevelWarn)", name), fn[0](common.LevelWarn), ansiYellow)
	colAssert(t, fmt.Sprintf("%s[0](common.LevelError)", name), fn[0](common.LevelError), ansiRed)
	colAssert(t, fmt.Sprintf("%s[0](common.LevelFatal)", name), fn[0](common.LevelFatal), ansiRed+ansiBold)
	colAssert(t, fmt.Sprintf("%s[1](common.LevelTrace)", name), fn[1](common.LevelTrace), ansiReset)
	colAssert(t, fmt.Sprintf("%s[1](common.LevelDebug)", name), fn[1](common.LevelDebug), ansiReset)
	colAssert(t, fmt.Sprintf("%s[1](common.LevelInfo)", name), fn[1](common.LevelInfo), ansiReset)
	colAssert(t, fmt.Sprintf("%s[1](common.LevelNotice)", name), fn[1](common.LevelNotice), ansiReset)
	colAssert(t, fmt.Sprintf("%s[1](common.LevelWarn)", name), fn[1](common.LevelWarn), ansiReset)
	colAssert(t, fmt.Sprintf("%s[1](common.LevelError)", name), fn[1](common.LevelError), ansiReset)
	colAssert(t, fmt.Sprintf("%s[1](common.LevelFatal)", name), fn[1](common.LevelFatal), ansiReset)
}

func TestColorOff(t *testing.T) {
	c := color.NewColor(common.ColorOff)
	colAssert(t, "c.TimFunc[0]()", c.TimFunc[0](), "")
	colAssert(t, "c.TimFunc[1]()", c.TimFunc[1](), "")
	colAssert(t, "c.UptFunc[0]()", c.UptFunc[0](), "")
	colAssert(t, "c.UptFunc[1]()", c.UptFunc[1](), "")
	colAssert(t, "c.LvlFunc[0]()", c.LvlFunc[0](), "")
	colAssert(t, "c.LvlFunc[1]()", c.LvlFunc[1](), "")
	colAssert(t, "c.SrcFunc[0]()", c.SrcFunc[0](), "")
	colAssert(t, "c.SrcFunc[1]()", c.SrcFunc[1](), "")
	colAssert(t, "c.MsgFunc[0]()", c.MsgFunc[0](), "")
	colAssert(t, "c.MsgFunc[1]()", c.MsgFunc[1](), "")
	colAssert(t, "c.KeyFunc[0]()", c.KeyFunc[0](), "")
	colAssert(t, "c.KeyFunc[1]()", c.KeyFunc[1](), "")
	colAssert(t, "c.ValFunc[0]()", c.ValFunc[0](), "")
	colAssert(t, "c.ValFunc[1]()", c.ValFunc[1](), "")
	colAssert(t, "c.ErKFunc[0]()", c.ErKFunc[0](), "")
	colAssert(t, "c.ErKFunc[1]()", c.ErKFunc[1](), "")
	colAssert(t, "c.ErVFunc[0]()", c.ErVFunc[0](), "")
	colAssert(t, "c.ErVFunc[1]()", c.ErVFunc[1](), "")
}

func TestColorSmart(t *testing.T) {
	c := color.NewColor(common.ColorSmart)
	colAssert(t, "c.TimFunc[0]()", c.TimFunc[0](), ansiWhite+ansiFaint)
	colAssert(t, "c.TimFunc[1]()", c.TimFunc[1](), ansiReset)
	colAssert(t, "c.UptFunc[0]()", c.UptFunc[0](), ansiCyan+ansiFaint)
	colAssert(t, "c.UptFunc[1]()", c.UptFunc[1](), ansiReset)
	dynAssert(t, "c.LvlFunc", c.LvlFunc)
	colAssert(t, "c.SrcFunc[0]()", c.SrcFunc[0](), ansiWhite+ansiFaint)
	colAssert(t, "c.SrcFunc[1]()", c.SrcFunc[1](), ansiReset)
	colAssert(t, "c.MsgFunc[0]()", c.MsgFunc[0](), ansiWhite)
	colAssert(t, "c.MsgFunc[1]()", c.MsgFunc[1](), ansiReset)
	colAssert(t, "c.KeyFunc[0]()", c.KeyFunc[0](), ansiWhite+ansiFaint)
	colAssert(t, "c.KeyFunc[1]()", c.KeyFunc[1](), ansiReset)
	colAssert(t, "c.ValFunc[0]()", c.ValFunc[0](), ansiWhite)
	colAssert(t, "c.ValFunc[1]()", c.ValFunc[1](), ansiReset)
	colAssert(t, "c.ErKFunc[0]()", c.ErKFunc[0](), ansiRed)
	colAssert(t, "c.ErKFunc[1]()", c.ErKFunc[1](), ansiReset)
	colAssert(t, "c.ErVFunc[0]()", c.ErVFunc[0](), ansiRed+ansiBold)
	colAssert(t, "c.ErVFunc[1]()", c.ErVFunc[1](), ansiReset)
}
