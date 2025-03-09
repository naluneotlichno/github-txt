package main

import (
	"fmt"
	"os"

	"githubtxt/file"
	"githubtxt/log"
	"githubtxt/repo"
	"githubtxt/utils"
)

func main() {
	mainTimer := utils.StartTimer()

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
	cloneTimer := utils.StartTimer()
	repo.CloneRepo(repoURL, repoPath, multiWriter)
	cloneTimer.PrintElapsedTime("клонирования", multiWriter)

	// Обработка файлов репозитория
	processTimer := utils.StartTimer()
	file.ProcessFiles(repoPath, savePath, multiWriter)
	processTimer.PrintElapsedTime("обработки файлов", multiWriter)

	// Удаление репозитория после обработки
	cleanupTimer := utils.StartTimer()
	repo.CleanupRepo(repoPath, multiWriter)
	cleanupTimer.PrintElapsedTime("удаления репозитория", multiWriter)

	mainTimer.PrintElapsedTime("всей программы", multiWriter)
}
