//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package example1

import (
	"time"

	// don't import anything else here
	"github.com/sfmunoz/logit"
)

func info() {
	log := logit.Logit()
	log.Info("log := logit.Logit().")
	log.Info("    WithLevel(logit.LevelNotice).")
	log.Info("    With(\"name\", \"example1\").")
	log.Info("    WithSymbolSet(logit.SymbolUnicodeDown).")
	log.Info("    WithUptimeFormat(logit.UptimeStd)")
}

func Run() {
	info()
	log := logit.Logit().
		WithLevel(logit.LevelNotice).
		With("name", "example1").
		WithSymbolSet(logit.SymbolUnicodeDown).
		WithUptimeFormat(logit.UptimeStd)
	log.Trace("trace-msg")
	log.Debug("debug-msg")
	log.Info("info-msg")
	log.Notice("notice-msg (trace, debug and info hidden because level=notice)", "duration", 2*time.Hour+5*time.Minute+3*time.Second+427*time.Millisecond)
	log.Warn("warn-msg")
	log.Error("error-msg")
}
