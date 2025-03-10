package repo

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// CloneRepo клонирует репозиторий
func CloneRepo(repoURL, repoPath string, log io.Writer) error {
	fmt.Fprintln(log, "🚀 Начинаем клонирование репозитория...")

	cmd := exec.Command("git", "clone", "--depth=1", "--filter=blob:none", repoURL, repoPath)
	cmd.Stdout = log
	cmd.Stderr = log
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(log, "Ошибка клонирования репозитория:", err)
		return err
	}
	return nil
}

// CleanupRepo удаляет клонированный репозиторий после обработки
func CleanupRepo(repoPath string, log io.Writer) error {
	if err := os.RemoveAll(repoPath); err != nil {
		fmt.Fprintln(log, "Ошибка удаления репозитория:", err)
		return err
	}
	fmt.Fprintln(log, "📂 Папка репозитория успешно удалена.")
	return nil
}
