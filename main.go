package main

import (
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

	var repoURL, savePath, repoPath string

	// 🔄 Повторяем ввод репозиторя в случае ошибки
	utils.HandleErrorRetry(func() error {
		var err error
		repoURL, savePath, repoPath, err = utils.InitRepo(multiWriter)
		return err
	}, "Ошибка инициализации репозитория", multiWriter, true)

	// 🔄 Клонирование репозитория
	utils.RunTimedAction(func() error {
		return repo.CloneRepo(repoURL, repoPath, multiWriter)
	}, "Клонирование репозитория", multiWriter, true)

	// 🔄 Обработка файлов
	utils.RunTimedAction(func() error {
		return file.ProcessFiles(repoPath, savePath, multiWriter)
	}, "Обработка файлов", multiWriter, true)

	// 🔄 Удаление репозитория
	utils.RunTimedAction(func() error {
		return repo.CleanupRepo(repoPath, multiWriter)
	}, "Удаление репозитория", multiWriter, true)

	mainTimer.PrintElapsedTime("всей программы", multiWriter)
}
