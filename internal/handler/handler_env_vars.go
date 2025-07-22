//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package handler

import (
	"io"
	"log/slog"
	"os"

	"github.com/sfmunoz/logit/internal/common"
)

func LogitWriterEnv() io.Writer {
	switch os.Getenv("LOGIT_WRITER") {
	case "stdout", "STDOUT":
		return os.Stdout
	case "stderr", "STDERR":
		return os.Stderr
	}
	return os.Stderr
}

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

func LogitUptimeFormatEnv() common.UptimeFormat {
	switch os.Getenv("LOGIT_UPTIME_FORMAT") {
	case "std", "STD":
		return common.UptimeStd
	case "adhoc", "ADHOC":
		return common.UptimeAdhoc
	}
	return common.UptimeAdhoc // default
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

func LogitSymbolSetEnv() common.SymbolSet {
	switch os.Getenv("LOGIT_SYMBOL_SET") {
	case "none", "NONE":
		return common.SymbolNone
	case "unicode_up", "UNICODE_UP":
		return common.SymbolUnicodeUp
	case "unicode_down", "UNICODE_DOWN":
		return common.SymbolUnicodeDown
	}
	return common.SymbolNone // default
}

func LogitAttrsModeEnv() common.AttrsMode {
	switch os.Getenv("LOGIT_ATTRS_MODE") {
	case "std", "STD":
		return common.AttrsStd
	case "builtin", "BUILTIN":
		return common.AttrsBuiltin
	}
	return common.AttrsStd // default
}
