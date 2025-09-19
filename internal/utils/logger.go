package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Logger struct {
	file *os.File
}

func NewLogger(filename string) *Logger {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	return &Logger{file: file}
}

func (l *Logger) LogRequest(method, endpoint, clientIP string, statusCode int, duration time.Duration) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] %s %s - %s - %d - %v\n", 
		timestamp, method, endpoint, clientIP, statusCode, duration)

	// Log to console
	fmt.Print(logMessage)

	// Log to file
	if l.file != nil {
		l.file.WriteString(logMessage)
	}
}

func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}
