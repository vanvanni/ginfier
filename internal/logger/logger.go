package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	ERROR LogLevel = "ERROR"
	FATAL LogLevel = "FATAL"
)

type LogEntry struct {
	Level     LogLevel `json:"level"`
	Timestamp string   `json:"timestamp"`
	Message   string   `json:"message"`
	File      string   `json:"file"`
}

func write(level LogLevel, message string) {
	_, filePath, line, _ := runtime.Caller(2)

	file := filepath.ToSlash(filePath)
	cwd, err := os.Getwd()
	if err == nil {
		cwd_cleaned := filepath.ToSlash(cwd)
		file = strings.TrimPrefix(file, cwd_cleaned)
	}

	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level,
		Message:   message,
		File:      fmt.Sprintf("%s:%d", file, line),
	}

	jsonEntry, err := json.Marshal(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling log entry: %v\n", err)
		return
	}

	os.Stdout.Write(jsonEntry)
	os.Stdout.Write([]byte("\n"))

	if level == FATAL {
		os.Exit(1)
	}
}

func Debug(message string) {
	write(DEBUG, message)
}

func Info(message string) {
	write(INFO, message)
}

func Error(message string) {
	write(ERROR, message)
}

func Fatal(message string) {
	write(FATAL, message)
	os.Exit(1)
}
