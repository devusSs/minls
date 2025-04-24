package cli

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/devusSs/minls/internal/log"
	"github.com/devusSs/minls/internal/storage"
)

func Clear() error {
	err := initialize()
	if err != nil {
		return fmt.Errorf("could not initialize: %w", err)
	}

	log.Debug("cli - Clear", slog.String("action", "initialized"))

	co := getClearOption()
	log.Debug(
		"cli - Clear",
		slog.String("action", "get_clear_option"),
		slog.String("co", co.String()),
	)

	err = handleClearOption(co)
	if err != nil {
		return fmt.Errorf("could not handle clear option: %w", err)
	}

	// don't log here to prevent issues
	// in case we delete all / logs

	return nil
}

type clearOption int

const (
	clearOptionInvalid clearOption = iota - 1
	clearOptionAll
	clearOptionData
	clearOptionLogs
	clearOptionDownloads
)

func (c clearOption) String() string {
	switch c {
	case clearOptionInvalid:
		return "invalid"
	case clearOptionAll:
		return "all"
	case clearOptionData:
		return "data"
	case clearOptionLogs:
		return "logs"
	case clearOptionDownloads:
		return "downloads"
	default:
		return "unknown"
	}
}

func getClearOption() clearOption {
	if len(os.Args) < 3 {
		return clearOptionInvalid
	}

	opt := strings.ToLower(os.Args[2])

	switch opt {
	case "all":
		return clearOptionAll
	case "data":
		return clearOptionData
	case "logs":
		return clearOptionLogs
	case "downloads":
		return clearOptionDownloads
	default:
		return clearOptionInvalid
	}
}

// TODO: implement downloads removal
func handleClearOption(opt clearOption) error {
	switch opt {
	case clearOptionInvalid:
		return errors.New("invalid clear option passed")
	case clearOptionAll:
		err := storage.RemoveStorageDir()
		if err != nil {
			return fmt.Errorf("could not remove storage dir: %w", err)
		}

		err = log.RemoveLogsDir()
		if err != nil {
			return fmt.Errorf("could not remove logs dir: %w", err)
		}

		return nil
	case clearOptionData:
		return storage.RemoveStorageDir()
	case clearOptionLogs:
		return log.RemoveLogsDir()
	case clearOptionDownloads:
		return errors.New("downloads not implemented yet")
	default:
		return fmt.Errorf("unknown clear option passed: %v", opt)
	}
}
