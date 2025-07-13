//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package logit_test

import (
	"testing"
	"time"

	// don't import anything else here -> black-box test
	"github.com/sfmunoz/logit"
)

func TestNewLogit(t *testing.T) {
	log := logit.Logit().
		WithLevel(logit.LevelNotice).
		With("test", "NewLogit").
		WithSymbolSet(logit.SymbolUnicodeDown).
		WithUptimeFormat(logit.UptimeStd)

	log.Trace("trace-msg")
	log.Debug("debug-msg")
	log.Info("info-msg")
	log.Notice("notice-msg", "duration", 2*time.Hour+5*time.Minute+3*time.Second+427*time.Millisecond)
	log.Warn("warn-msg")
	log.Error("error-msg")
}
