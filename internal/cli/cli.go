package cli

import (
	"os"

	"github.com/devusSs/minls/internal/log"
)

func initialize() error {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	err := log.SetupLog(logLevel)
	if err != nil {
		return err
	}

	return nil
}
