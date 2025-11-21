package utils

import (
	"os"

	"github.com/UncleJunVIP/gabagool/pkg/gabagool"
	"github.com/UncleJunVIP/gabagool/pkg/gabagool/i18n"
)

func init() {
	err := i18n.InitI18N([]string{"resources/i18n/en.json", "resources/i18n/es.json"})
	if err != nil {
		gabagool.GetLogger().Error("Failed to initialize i18n", "error", err)
		os.Exit(1)
	}
}
