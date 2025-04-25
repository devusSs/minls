package clip

import (
	"fmt"
	"log/slog"
	"runtime"

	"github.com/atotto/clipboard"

	"github.com/devusSs/minls/internal/log"
)

func Init() error {
	log.Debug("clip - Init", slog.String("go_os", runtime.GOOS))

	switch runtime.GOOS {
	case "windows":
	case "darwin":
	case "linux":
	default:
		return fmt.Errorf("clip not implemented for %s", runtime.GOOS)
	}

	return nil
}

func Read() (string, error) {
	log.Debug("clip - Read", slog.String("go_os", runtime.GOOS))
	return clipboard.ReadAll()
}

func Write(input string) error {
	log.Debug("clip - Write", slog.String("go_os", runtime.GOOS), slog.String("input", input))
	return clipboard.WriteAll(input)
}
