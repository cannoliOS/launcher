package utils

import (
	"cannoliOS/models"
	"cannoliOS/state"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/UncleJunVIP/gabagool/pkg/gabagool"
	"github.com/UncleJunVIP/gabagool/pkg/gabagool/i18n"
)

func LoadConfig() error {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return err
	}

	var config models.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	gabagool.SetLogLevel(config.LogLevel)
	i18n.SetWithCode(config.Language)
	state.Init(&config)

	return nil
}

func GetLoggerInstance() *slog.Logger {
	return gabagool.GetLoggerInstance()
}

func GetConfig() *models.Config {
	return state.Get().Config
}
