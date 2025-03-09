package utils

import (
	"bufio"         // чтение файлов и ввода потрочно
	"fmt"
	"io"            // для работы с потоками данных
	"os"            // для работы с операционной системой
	"path/filepath" // для работы с путями
	"strings"
)

// GetRepoURL запрашивает у пользователя URL репозитория
func GetRepoURL(log io.Writer) string {
	fmt.Fprint(log, "\nВведите URL GitHub-репозитория (или нажмите Enter для выхода): ")
	reader := bufio.NewReader(os.Stdin)
	repoURL, _ := reader.ReadString('\n')
	return strings.TrimSpace(repoURL)
}

// GetRepoNameFromURL получает имя репозитория
func GetRepoNameFromURL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) > 1 {
		return strings.TrimSuffix(parts[len(parts)-1], ".git")
	}
	return ""
}

// GetSavePath возвращает путь к папке, в которую сохраняются файлы
func GetSavePath(repoName string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Ошибка получения домашней директории:", err)
		os.Exit(1)
	}
	return filepath.Join(homeDir, "Desktop", repoName)
}

// GetRepoPath возвращает путь к клонированному репозиторию
func GetRepoPath(savePath string) string {
	return filepath.Join(savePath, "repo")
}
