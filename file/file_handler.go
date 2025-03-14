package file

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"unicode/utf8"
)

// ProcessFiles записывает все файлы с текском(txt) в один txt файл
func ProcessFiles(repoPath, outputFilePath string, log io.Writer) error {
	fmt.Fprintln(log, "📂 Обрабатываем файлы по пакетам...")

	// Создаем файл, в который будем записывать все данные
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Fprintln(log, "❌ Ошибка создания выходного файла:", err)
		return err
	}
	defer outputFile.Close()

	// Делаем обход всех файлов в репозитории через рекурсивную функцию
	err = filepath.WalkDir(repoPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprint(log, "❌ Ошибка обхода файлов:", err)
			return err
		}

		// Пропускаем мусорные папки
		if d.IsDir() {
			skipFolder := map[string]bool{
				".git":    true,
				".vscode": true,
				".idea":   true,
				".cache":  true,
				".env":    true,
			}
			if skipFolder[filepath.Base(path)] {
				fmt.Fprintln(log, "📂 Пропускаем папку:", path)
				return filepath.SkipDir
			}
		}

		// Обрабатываем файл если он не папка
		if !d.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				fmt.Fprintln(log, "❌ Ошибка чтения файла:", path, err)
				return nil
			}

			// Пропускаем файл если он не UTF-8
			if !utf8.Valid(content) {
				fmt.Fprintln(log, "📖 Файл не содержит кодировки UTF-8:", path)
				return nil
			}

			// Записываем содержимое файла в общий TXT файл
			_, err = outputFile.WriteString(fmt.Sprintf("\nFile: %s\n\n%s\n", path, string(content)))
			if err != nil {
				fmt.Fprintln(log, "❌ Ошибка записи файла:", outputFilePath, err)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Fprintln(log, "⚠️ Обнаружены ошибки при обработке файлов:", err)
		return err
	}

	fmt.Fprintln(log, "✅ Обработка файлов завершена!")
	return nil
}
