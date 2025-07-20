//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package example3

import (
	"context"
	"log/slog"

	// don't import anything else here
	"github.com/sfmunoz/logit"
)

func info() {
	log := logit.Logit()
	log.Info("log := logit.Logit().")
	log.Info("    With(\"k1\", \"v1\").")
	log.Info("    WithGroup(\"g1\").")
	log.Info("    With(\"k2\", \"v2\").")
	log.Info("    WithGroup(\"g2\").")
	log.Info("    WithGroup(\"g3\")")
}

func Run() {
	info()
	log := logit.Logit().
		With("k1", "v1").
		WithGroup("g1").
		With("k2", "v2").
		WithGroup("g2").
		WithGroup("g3")
	log.
		WithGroup("g11").
		LogAttrs(
			context.Background(),
			logit.LevelNotice, "(g11)",
			slog.Int("a", 1),
			slog.Int("b", 2),
		)
	log.
		LogAttrs(
			context.Background(),
			logit.LevelNotice,
			"(g12)",
			slog.Group(
				"g12",
				slog.Int("a", 1),
				slog.Int("b", 2),
			),
		)
}
