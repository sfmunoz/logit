//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package main

import (
	"reflect"
	"runtime"
	"time"

	// don't import anything else here
	"github.com/sfmunoz/logit"
)

var log = logit.Logit()

func funcName(fn any) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

func example1() {
	log.Info("log := logit.Logit().")
	log.Info("    WithLevel(logit.LevelNotice).")
	log.Info("    With(\"name\", \"example1\").")
	log.Info("    WithSymbolSet(logit.SymbolUnicodeDown).")
	log.Info("    WithUptimeFormat(logit.UptimeStd)")
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

func example2() {
	log.Info("logit.Logit().")
	log.Info("WithLevel(logit.LevelTrace).")
	log.Info("WithSymbolSet(logit.SymbolUnicodeDown).")
	log.Info("WithTpl(logit.TplUptime, logit.TplLevel, logit.TplSource)")
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

func main() {
	examples := []func(){example1, example2}
	for _, f := range examples {
		fName := funcName(f)
		log.Info("================ " + fName + " ================")
		f()
		log.Info("---------------- " + fName + " ----------------")
	}
}
