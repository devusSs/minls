package log

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Init() error {
	err := createLogsDirIfNotExist()
	if err != nil {
		return fmt.Errorf("could not create logs dir: %w", err)
	}

	go cleanOldLogs()

	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}

	err = createFileLogger(level)
	if err != nil {
		return fmt.Errorf("could not create file logger: %w", err)
	}

	createConsoleLogger(level)

	setup = true

	if level == "debug" {
		Warn("log - Init", slog.String("warn", "debug logging enabled, might leak sensitive data"))
	}

	return nil
}

func Debug(msg string, args ...any) {
	if !setup {
		fmt.Println("log - Debug: not setup, cannot log")
		return
	}

	fileLogger.Debug(msg, args...)
	consoleLogger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	if !setup {
		fmt.Println("logs - Info: not setup, cannot log")
		return
	}

	fileLogger.Info(msg, args...)
	consoleLogger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	if !setup {
		fmt.Println("logs - Warn: not setup, cannot log")
		return
	}

	fileLogger.Warn(msg, args...)
	consoleLogger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	if !setup {
		fmt.Println("logs - Error: not setup, cannot log")
		return
	}

	fileLogger.Error(msg, args...)
	consoleLogger.Error(msg, args...)
}

var setup bool

var logsDir = "logs"

func createLogsDirIfNotExist() error {
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not find executeable: %w", err)
	}

	logsDir = filepath.Join(filepath.Dir(exe), logsDir)

	err = os.Mkdir(logsDir, 0755)
	if err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("could not create logs dir: %w", err)
		}
	}

	return nil
}

var fileLogger *slog.Logger

func createFileLogger(level string) error {
	logFileName := "minls_" + time.Now().Format("2006-01-02_15-04-05") + ".log.json"
	logFilePath := filepath.Join(logsDir, logFileName)

	logFile, err := os.Create(logFilePath)
	if err != nil {
		return fmt.Errorf("could not create log file: %w", err)
	}

	fileLogger = slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level: slogLevelFromString(level),
	}))

	return nil
}

var consoleLogger *slog.Logger

func createConsoleLogger(level string) {
	consoleLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slogLevelFromString(level),
	}))
}

func slogLevelFromString(s string) slog.Level {
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
		fmt.Printf("log: unknown log level provided '%s', using default 'info'\n", s)
		return slog.LevelInfo
	}
}

const logFileMaxAge = 7 * 24 * time.Hour

// TODO: find a way to complete this before program exits on main
func cleanOldLogs() {
	entries, err := os.ReadDir(logsDir)
	if err != nil {
		fmt.Println("log - cleanOldLogs: could not read dir:", err)
		return
	}

	for _, entry := range entries {
		fp := filepath.Join(logsDir, entry.Name())

		if entry.IsDir() {
			fmt.Println(
				"log - cleanOldLogs: found unexpected directory:",
				fp,
			)
			continue
		}

		split := strings.SplitN(fp, ".", 2)
		if len(split) != 2 {
			fmt.Printf(
				"log - cleanOldLogs: found malformed file: '%s', len(split) = %d\n",
				fp,
				len(split),
			)
			continue
		}

		if split[1] != "log.json" {
			fmt.Printf("log - cleanOldLogs: found unexpected file: '%s'\n", fp)
			continue
		}

		fi, err := os.Stat(filepath.Join(fp))
		if err != nil {
			fmt.Printf(
				"log - cleanOldLogs: could not get file info for file '%s', error: %v\n",
				fp,
				err,
			)
			continue
		}

		if time.Since(fi.ModTime()) > logFileMaxAge {
			err = os.Remove(fp)
			if err != nil {
				fmt.Printf("log - cleanOldLogs: could not delete file: '%s', error: %v\n", fp, err)
				continue
			}
		}
	}
}
