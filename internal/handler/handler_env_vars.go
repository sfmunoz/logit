//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package handler

import (
	"os"

	"github.com/sfmunoz/logit/internal/common"
)

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
