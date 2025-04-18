package cli

import (
	"fmt"

	"github.com/devusSs/minls/internal/log"
)

func initialize() error {
	err := log.Init()
	if err != nil {
		return fmt.Errorf("could not init log: %w", err)
	}

	return nil
}
