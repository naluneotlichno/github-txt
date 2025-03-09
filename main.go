package main

import (
	"fmt"
	"githubtxt/file"
	"githubtxt/log"
	"githubtxt/repo"
	"githubtxt/utils"
	"os"
)

func main() {
	startTime := utils.StartTimer()

	// Создаём лог-файл
	logFile := log.CreateLogFile()
	defer logFile.Close()

	multiWriter := log.SetupLogger(logFile)

	repoURL := utils.GetRepoURL(multiWriter)
	if repoURL == "" {
		fmt.Fprintln(multiWriter, "👋 Завершение работы программы.")
		return
	}

	repoName := utils.GetRepoNameFromURL(repoURL)
	savePath := utils.GetSavePath(repoName)
	repoPath := utils.GetRepoPath(savePath)

	if err := os.RemoveAll(repoPath); err != nil {
		fmt.Fprintln(multiWriter, "Ошибка удаления папки репозитория:", err)
		return
	}

	// Клонирование репозитория
	repo.CloneRepo(repoURL, repoPath, multiWriter)

	// Обработка файлов
	file.ProcessFiles(repoPath, savePath, multiWriter)

	// Удаление репозитория после обработки
	repo.CleanupRepo(repoPath, multiWriter)

	utils.PrintExecutionTime(startTime, multiWriter)
}
