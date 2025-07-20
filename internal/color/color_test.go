//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//
//

package color_test

import (
	"testing"

	"github.com/sfmunoz/logit/internal/color"
	"github.com/sfmunoz/logit/internal/common"
)

func testCol(t *testing.T, msg, got, want string) {
	if got == want {
		return
	}
	t.Fatalf("%s failed: want='%s', got='%s'", msg, want, got)
}

func TestColorOff(t *testing.T) {
	c := color.NewColor(common.ColorOff)
	testCol(t, "c.TimFunc[0]()", c.TimFunc[0](), "")
	testCol(t, "c.TimFunc[1]()", c.TimFunc[1](), "")
	testCol(t, "c.UptFunc[0]()", c.UptFunc[0](), "")
	testCol(t, "c.UptFunc[1]()", c.UptFunc[1](), "")
	testCol(t, "c.LvlFunc[0]()", c.LvlFunc[0](), "")
	testCol(t, "c.LvlFunc[1]()", c.LvlFunc[1](), "")
	testCol(t, "c.SrcFunc[0]()", c.SrcFunc[0](), "")
	testCol(t, "c.SrcFunc[1]()", c.SrcFunc[1](), "")
	testCol(t, "c.MsgFunc[0]()", c.MsgFunc[0](), "")
	testCol(t, "c.MsgFunc[1]()", c.MsgFunc[1](), "")
	testCol(t, "c.KeyFunc[0]()", c.KeyFunc[0](), "")
	testCol(t, "c.KeyFunc[1]()", c.KeyFunc[1](), "")
	testCol(t, "c.ValFunc[0]()", c.ValFunc[0](), "")
	testCol(t, "c.ValFunc[1]()", c.ValFunc[1](), "")
	testCol(t, "c.ErKFunc[0]()", c.ErKFunc[0](), "")
	testCol(t, "c.ErKFunc[1]()", c.ErKFunc[1](), "")
	testCol(t, "c.ErVFunc[0]()", c.ErVFunc[0](), "")
	testCol(t, "c.ErVFunc[1]()", c.ErVFunc[1](), "")
}
