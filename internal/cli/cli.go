package cli

import (
	"fmt"
	"log/slog"

	"github.com/devusSs/minls/internal/clip"
	"github.com/devusSs/minls/internal/log"
	"github.com/devusSs/minls/internal/storage"
)

func initialize() error {
	err := log.Init()
	if err != nil {
		return fmt.Errorf("could not init log: %w", err)
	}

	log.Debug("cli - initialize", slog.String("action", "log_init"))

	err = storage.Init()
	if err != nil {
		return fmt.Errorf("could not init storage: %w", err)
	}

	log.Debug("cli - initialize", slog.String("action", "storage_init"))

	err = clip.Init()
	if err != nil {
		return fmt.Errorf("could not init clip: %w", err)
	}

	log.Debug("cli - initialize", slog.String("action", "clip_init"))

	return nil
}
