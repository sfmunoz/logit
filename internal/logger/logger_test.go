//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package logger_test

import (
	"context"
	//"errors"
	"fmt"
	//"log"
	"log/slog"
	"testing"
	//"time"

	"github.com/sfmunoz/logit/internal/common"
	"github.com/sfmunoz/logit/internal/logger"
)

func replaceAttr(_ []string, a slog.Attr) slog.Attr {
	if a.Key == "the-key" {
		return slog.String("the-key", fmt.Sprintf("'%s'=>EXTENDED BY 'replaceAttr()'", a.Value))
	}
	return slog.Attr{}
}

func inner(_ *testing.T, symbolSet common.SymbolSet) {
	l := logger.NewLogger(nil).
		With("a1", "v1").
		WithGroup("g1").
		WithGroup("g2").
		WithLevel(common.LevelTrace).
		WithSymbolSet(symbolSet).
		WithReplaceAttr(replaceAttr)
	//slog.SetDefault(l.Logger)
	//l.Info("symbols", "SymbolNone", common.SymbolNone, "SymbolUnicodeUp", common.SymbolUnicodeUp, "SymbolUnicodeDown", common.SymbolUnicodeDown, "Current", symbolSet)
	//l.Info("logger.NewLogger()", "type", fmt.Sprintf("%T", l))
	//slog.Info("Starting server", "addr", ":8080", "env", "production")
	//slog.Debug("Connected to DB", "db", "myapp", "host", "localhost:5432")
	//slog.Warn("Slow request", "method", "GET", "path", "/users", "duration", 497*time.Millisecond)
	//slog.Error("DB connection lost", "err", "connection reset", "failure", errors.New("network off"), "db", "myapp")
	//log.Print("log.Print() message")
	//l.Trace("trace", "the-key", "the-val")
	//l.WithTpl(common.TplLevel, common.TplUptime, common.TplUptime, common.TplLevel).Notice("notice (ad hoc template)", "the-key", "the-val")
	//l.Fatal("fatal", "key", "val")
	fmt.Println("---- 0010 ----")
	l.WithGroup("s1").LogAttrs(context.Background(), common.LevelNotice, "(1) logger.WithGroup(\"s1\")", slog.Int("a", 1), slog.Int("b", 2))
	fmt.Println("---- 0020 ----")
	l.LogAttrs(context.Background(), common.LevelNotice, "(2) logger.WithGroup(\"s2\")", slog.Group("s2", slog.Int("a", 1), slog.Int("b", 2)))
	fmt.Println("---- 0030 ----")
	//slog.Log(context.Background(), common.LevelNotice, "slog.Log(LevelNotice)", "the-key", "the-val")
}

func TestPlain(t *testing.T) {
	inner(t, common.SymbolNone)
}

func TestUnicodeUp(t *testing.T) {
	inner(t, common.SymbolUnicodeUp)
}

func TestUnicodeDown(t *testing.T) {
	inner(t, common.SymbolUnicodeDown)
}
