//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package handler

import (
	"log/slog"
	"os"

	"github.com/sfmunoz/logit/internal/common"
)

func LogitLevelEnv() slog.Level {
	switch os.Getenv("LOGIT_LEVEL") {
	case "trace", "TRACE":
		return common.LevelTrace
	case "debug", "DEBUG":
		return common.LevelDebug
	case "info", "INFO":
		return common.LevelInfo
	case "notice", "NOTICE":
		return common.LevelNotice
	case "warn", "WARN":
		return common.LevelWarn
	case "error", "ERROR":
		return common.LevelError
	case "fatal", "FATAL":
		return common.LevelFatal
	}
	return common.LevelInfo // default
}

func LogitTimeFormatEnv() string {
	ret := os.Getenv("LOGIT_TIME_FORMAT")
	if ret != "" {
		return ret
	}
	return "2006-01-02T15:04:05.000Z07:00"
}

func LogitColorModeEnv() common.ColorMode {
	switch os.Getenv("LOGIT_COLOR_MODE") {
	case "off", "OFF":
		return common.ColorOff
	case "smart", "SMART":
		return common.ColorSmart
	case "medium", "MEDIUM":
		return common.ColorMedium
	case "full", "FULL":
		return common.ColorFull
	}
	return common.ColorSmart // default
}
