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

// ProcessFiles обрабатывает файлы по пакетам
func ProcessFiles(repoPath, outputDir string, log io.Writer) {
	fmt.Fprintln(log, "📂 Обрабатываем файлы по пакетам...")

	packages := make(map[string]*os.File)
	var mu sync.Mutex
	var wg sync.WaitGroup
	files := make(chan string, 100)

	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("❌ Ошибка обхода файлов:", err)
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
		fmt.Fprintln(log, "Ошибка обхода файлов:", err)
		return
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range files {
				content, err := os.ReadFile(path)
				if err != nil {
					fmt.Fprintln(log, "❌ Ошибка чтения файла:", path, err)
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
						mu.Unlock()
						continue
					}
					packages[packageName] = f
				}
				f := packages[packageName]
				mu.Unlock()

				mu.Lock()
				_, _ = f.WriteString(fmt.Sprintf("\nFile: %s\n\n%s\n", path, string(content)))
				mu.Unlock()
			}
		}()
	}

	close(files)
	wg.Wait()

	for _, f := range packages {
		_ = f.Close()
	}

	fmt.Fprintln(log, "✅ Обработка файлов завершена!")
}
