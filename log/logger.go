package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// CreateLogFile создаёт лог-файл
func CreateLogFile() *os.File {
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)

	logsPath := filepath.Join(exeDir, "logs")
	_ = os.MkdirAll(logsPath, os.ModePerm)

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFilePath := filepath.Join(logsPath, timestamp+".log")

	logFile, err := os.Create(logFilePath)
	if err != nil {
		fmt.Println("❌ Ошибка создания лог-файла:", err)
		os.Exit(1)
	}

	return logFile
}

// SetupLogger настраивает логирование
func SetupLogger(logFile *os.File) io.Writer {
	return io.MultiWriter(os.Stdout, logFile)
}
