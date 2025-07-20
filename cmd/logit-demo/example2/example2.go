//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package example2

import (
	// don't import anything else here
	"github.com/sfmunoz/logit"
)

func info() {
	log := logit.Logit()
	log.Info("logit.Logit().")
	log.Info("    WithLevel(logit.LevelTrace).")
	log.Info("    WithSymbolSet(logit.SymbolUnicodeDown).")
	log.Info("    WithTpl(logit.TplUptime, logit.TplLevel, logit.TplSource)")
}

func Run() {
	info()
	log := logit.Logit().
		WithLevel(logit.LevelTrace).
		WithSymbolSet(logit.SymbolUnicodeDown).
		WithTpl(logit.TplUptime, logit.TplLevel, logit.TplSource)
	log.Info("every source ref must be OK (both 'logit' and 'slog')")
	log.Trace("trace (logit)")
	log.Debug("debug (slog)")
	log.Info("info (slog)")
	log.Notice("notice (logit)")
	log.Warn("warn (slog)")
	log.Error("error (slog) - 'fatal' and 'panic' skipped to prevent program from finishing")
	//log.Fatal("fatal (logit)")
	//log.Panic("panic (logit)")
}
