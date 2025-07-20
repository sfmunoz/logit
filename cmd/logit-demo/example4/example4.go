//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package example4

import (
	"log/slog"
	"os"

	// don't import anything else here
	"github.com/sfmunoz/logit"
)

func info() {
	log := logit.Logit()
	log.Info("log := logit.Logit().WithHandlers(")
	log.Info("    []slog.Handler{")
	log.Info("        slog.NewTextHandler(")
	log.Info("            os.Stderr,")
	log.Info("            &slog.HandlerOptions{AddSource: true, Level: logit.LevelInfo},")
	log.Info("        ),")
	log.Info("        slog.NewJSONHandler(")
	log.Info("            os.Stderr,")
	log.Info("            &slog.HandlerOptions{AddSource: false, Level: logit.LevelInfo},")
	log.Info("        ),")
	log.Info("    },")
	log.Info(")")
}

func Run() {
	info()
	log := logit.Logit().WithHandlers(
		[]slog.Handler{
			slog.NewTextHandler(
				os.Stderr,
				&slog.HandlerOptions{AddSource: true, Level: logit.LevelInfo},
			),
			slog.NewJSONHandler(
				os.Stderr,
				&slog.HandlerOptions{AddSource: false, Level: logit.LevelInfo},
			),
		},
	)
	log.Info("Message repeated", "times", 3)
}
