//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package logit_test

import (
	"testing"

	// don't import anything else here -> black-box test
	"github.com/sfmunoz/logit"
)

func TestNewLogit(t *testing.T) {
	log := logit.Logit().
		WithLevel(logit.LevelNotice).
		With("test", "NewLogit")
	log.Trace("trace-msg")
	log.Debug("debug-msg")
	log.Info("info-msg")
	log.Notice("notice-msg")
	log.Warn("warn-msg")
	log.Error("error-msg")
}
