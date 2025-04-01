package log

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SetupLog(level string) error {
	err := createLogsDir()
	if err != nil {
		return fmt.Errorf("could not create logs dir: %w", err)
	}

	go removeOldLogs()

	createConsoleLogger(level)

	err = createFileLogger(level)
	if err != nil {
		return fmt.Errorf("could not create file logger: %w", err)
	}

	setup = true

	return nil
}

func Debug(msg string, args ...any) {
	if !setup {
		fmt.Println("log: not setup, cannot log debug")
		return
	}

	consoleLogger.Debug(msg, args...)
	fileLogger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	if !setup {
		fmt.Println("log: not setup, cannot log info")
		return
	}

	consoleLogger.Info(msg, args...)
	fileLogger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	if !setup {
		fmt.Println("log: not setup, cannot log warn")
		return
	}

	consoleLogger.Warn(msg, args...)
	fileLogger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	if !setup {
		fmt.Println("log: not setup, cannot log error")
		return
	}

	consoleLogger.Error(msg, args...)
	fileLogger.Error(msg, args...)
}

var (
	consoleLogger *slog.Logger
	fileLogger    *slog.Logger
)

var setup bool

var logsDir = "logs"

func createLogsDir() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not find executable path: %w", err)
	}

	logsDir = filepath.Join(filepath.Dir(exe), logsDir)

	err = os.Mkdir(logsDir, 0750)
	if err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("could not create logs dir: %w", err)
		}
	}

	return nil
}

const logFileMaxAge = 30 * 24 * time.Hour

func removeOldLogs() {
	files, err := os.ReadDir(logsDir)
	if err != nil {
		fmt.Println("log: could not read logs dir:", err)
		return
	}

	for _, file := range files {
		var fi fs.FileInfo
		fi, err = file.Info()
		if err != nil {
			fmt.Println("log: could not read file info:", err)
			continue
		}

		if time.Since(fi.ModTime()) > logFileMaxAge {
			err = os.Remove(filepath.Join(logsDir, fi.Name()))
			if err != nil {
				fmt.Println("log: could not delete file:", err)
				continue
			}
		}
	}
}

func createConsoleLogger(level string) {
	consoleLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseSlogLevelFromString(level),
	}))
}

const logFormat = "2006-01-02_15-04-05"

var logFile = fmt.Sprintf("minls_%s.json.log", time.Now().Format(logFormat))

func createFileLogger(level string) error {
	f, err := os.Create(filepath.Join(logsDir, logFile))
	if err != nil {
		return fmt.Errorf("could not create log file: %w", err)
	}

	fileLogger = slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{
		Level: parseSlogLevelFromString(level),
	}))

	return nil
}

var defaultLogLevel = slog.LevelInfo

func parseSlogLevelFromString(s string) slog.Level {
	s = strings.ToLower(s)

	switch s {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		fmt.Printf("log: invalid log level passed: %s, using info\n", s)
		return defaultLogLevel
	}
}
