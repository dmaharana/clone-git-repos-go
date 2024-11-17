package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Standard log flags
const (
	Ldate         = log.Ldate         // the date in the local time zone: 2009/01/23
	Ltime         = log.Ltime         // the time in the local time zone: 01:23:23
	Lmicroseconds = log.Lmicroseconds // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile     = log.Llongfile     // full file name and line number: /a/b/c/d.go:23
	Lshortfile    = log.Lshortfile    // final file name element and line number: d.go:23
	LUTC          = log.LUTC          // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags     = log.LstdFlags     // initial values for the standard logger
)

// Logger wraps the standard log package
type Logger struct {
	*log.Logger
	file *os.File
}

// Config holds logger configuration
type Config struct {
	LogDir     string
	MaxSize    int64  // in bytes
	TimeFormat string
}

// DefaultConfig returns default logger configuration
func DefaultConfig() *Config {
	return &Config{
		LogDir:     "logs",
		MaxSize:    10 * 1024 * 1024, // 10MB
		TimeFormat: "2006-01-02",
	}
}

// New creates a new logger instance
func New(cfg *Config) (*Logger, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	// Generate log file name with current date
	logFileName := fmt.Sprintf("clone-git-repo-%s.log", time.Now().Format(cfg.TimeFormat))
	logFilePath := filepath.Join(cfg.LogDir, logFileName)

	// Open log file
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	// Create multi-writer to write to both file and stdout
	multiWriter := os.Stdout
	if file != nil {
		multiWriter = file
	}

	// Create logger
	logger := &Logger{
		Logger: log.New(multiWriter, "", log.LstdFlags|log.Lshortfile),
		file:   file,
	}

	return logger, nil
}

// Close closes the log file
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// Printf formats and prints a message to the log
func (l *Logger) Printf(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...))
}

// Print prints a message to the log
func (l *Logger) Print(v ...interface{}) {
	l.Output(2, fmt.Sprint(v...))
}

// Println prints a message to the log with a newline
func (l *Logger) Println(v ...interface{}) {
	l.Output(2, fmt.Sprintln(v...))
}

// Fatal prints a message and calls os.Exit(1)
func (l *Logger) Fatal(v ...interface{}) {
	l.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf formats and prints a message and calls os.Exit(1)
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}
