//
// Author: 46285520+sfmunoz@users.noreply.github.com
// URL:    https://github.com/sfmunoz/logit
//

package handler

import (
	"os"

	"github.com/sfmunoz/logit/internal/common"
)

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
