package utils

import (
	"bufio"
	"fmt"
	"io"            
	"os"            
	"path/filepath" 
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

// InitRepo выполняет начальную настройку репозитория
func InitRepo(log io.Writer) (repoURL, savePath, repoPath string, err error) {
	repoURL = GetRepoURL(log)
	if repoURL == "" {
		fmt.Fprintln(log, "ВведитеURL GitHub-репозитория.")
		return "", "", "", nil
	}

	repoName := GetRepoNameFromURL(repoURL)
	savePath = GetSavePath(repoName)
	repoPath = GetRepoPath(savePath)

	if err := os.RemoveAll(repoPath); err != nil {
		return "", "", "", fmt.Errorf("не удалось удалить папку репозиторияЖ %w", err)
	}

	return repoURL, savePath, repoPath, nil
}

// HandleErrorRetry обрабатывает ошибку и повторяет попытку повторно
func HandleErrorRetry(action func() error, msg string, log io.Writer, retry bool) {
	for {
		err := action()
		if err == nil {
			break
		}

		fmt.Fprintln(log, "❌", msg, ":", err)

		if retry {
			fmt.Fprintln(log, "🔄 Повторяем попытку...")
			continue
		} else {
			fmt.Fprintln(log, "💥 Критическая ошибка. Завершение работы.")
			os.Exit(1)
		}
	}
}
