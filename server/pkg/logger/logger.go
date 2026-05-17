// Package logger 提供結構化日誌
package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Level 日誌等級
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var currentLevel = INFO

// SetLevel 設定日誌等級
func SetLevel(l Level) {
	currentLevel = l
}

func logf(level Level, prefix, format string, args ...interface{}) {
	if level < currentLevel {
		return
	}
	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	log.Printf("[%s] %s %s", timestamp, prefix, msg)
}

// Debug 除錯日誌
func Debug(format string, args ...interface{}) {
	logf(DEBUG, "DEBUG", format, args...)
}

// Info 資訊日誌
func Info(format string, args ...interface{}) {
	logf(INFO, "INFO ", format, args...)
}

// Warn 警告日誌
func Warn(format string, args ...interface{}) {
	logf(WARN, "WARN ", format, args...)
}

// Error 錯誤日誌
func Error(format string, args ...interface{}) {
	logf(ERROR, "ERROR", format, args...)
}

// Fatal 致命錯誤（會結束程式）
func Fatal(format string, args ...interface{}) {
	logf(ERROR, "FATAL", format, args...)
	os.Exit(1)
}
