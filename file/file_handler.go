package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode/utf8"
)

// ProcessFiles обрабатывает файлы по пакетам и возвращает ошибку при сбоях
func ProcessFiles(repoPath, outputDir string, log io.Writer) error {
	fmt.Fprintln(log, "📂 Обрабатываем файлы по пакетам...")

	packages := make(map[string]*os.File)
	var mu sync.Mutex
	var wg sync.WaitGroup
	files := make(chan string, 100)
	errors := make(chan error, 10) // Канал для ошибок

	// Обход файлов
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintln(log, "❌ Ошибка обхода файлов:", err)
			return err
		}
		if strings.Contains(path, ".git") {
			return filepath.SkipDir
		}
		if !info.IsDir() {
			files <- path
		}
		return nil
	})

	if err != nil {
		fmt.Fprintln(log, "❌ Ошибка обхода файлов:", err)
		return err
	}

	// Воркеры для обработки файлов
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range files {
				content, err := os.ReadFile(path)
				if err != nil {
					fmt.Fprintln(log, "❌ Ошибка чтения файла:", path, err)
					errors <- fmt.Errorf("ошибка чтения файла %s: %w", path, err)
					continue
				}

				if !utf8.Valid(content) {
					fmt.Fprintln(log, "🚫 Пропущен бинарный файл:", path)
					continue
				}

				packageName := filepath.Base(filepath.Dir(path))
				outputFilePath := filepath.Join(outputDir, "repo_"+packageName+".txt")

				mu.Lock()
				if _, exists := packages[packageName]; !exists {
					f, err := os.Create(outputFilePath)
					if err != nil {
						fmt.Fprintln(log, "❌ Ошибка создания файла пакета:", outputFilePath, err)
						errors <- fmt.Errorf("ошибка создания файла %s: %w", outputFilePath, err)
						mu.Unlock()
						continue
					}
					packages[packageName] = f
				}
				f := packages[packageName]
				mu.Unlock()

				mu.Lock()
				_, err = f.WriteString(fmt.Sprintf("\nFile: %s\n\n%s\n", path, string(content)))
				if err != nil {
					fmt.Fprintln(log, "❌ Ошибка записи в файл:", outputFilePath, err)
					errors <- fmt.Errorf("ошибка записи в файл %s: %w", outputFilePath, err)
				}
				mu.Unlock()
			}
		}()
	}

	close(files)
	wg.Wait()
	close(errors) // Закрываем канал ошибок после завершения всех воркеров

	// Проверяем, были ли ошибки в обработке файлов
	var finalErr error
	for err := range errors {
		if finalErr == nil {
			finalErr = err // Запоминаем первую ошибку
		} else {
			finalErr = fmt.Errorf("%v; %w", finalErr, err) // Объединяем ошибки
		}
	}

	// Закрываем все файлы
	for _, f := range packages {
		_ = f.Close()
	}

	if finalErr != nil {
		fmt.Fprintln(log, "⚠️ Обнаружены ошибки при обработке файлов:", finalErr)
		return finalErr
	}

	fmt.Fprintln(log, "✅ Обработка файлов завершена!")
	return nil
}
