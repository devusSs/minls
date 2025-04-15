package clip

import (
	"fmt"
	"runtime"

	"github.com/atotto/clipboard"
)

func Init() error {
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
	return clipboard.ReadAll()
}

func Write(input string) error {
	return clipboard.WriteAll(input)
}
