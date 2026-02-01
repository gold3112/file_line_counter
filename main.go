package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Configuration for skipping files and directories
var (
	skipDirs = []string{".git", "node_modules", "vendor", "build", "dist"}
	skipFiles = []string{"region_database.json"}
	codeExtensions = []string{
		".go", ".py", ".java", ".c", ".cpp", ".h", ".hpp", ".cs", ".js", ".ts",
		".html", ".css", ".scss", ".rb", ".php", ".swift", ".kt", ".rs", ".sh",
	}
)

// Command-line flags
var (
	codeOnly = flag.Bool("codeonly", false, "指定された拡張子のソースコードファイルのみをカウントします。")
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("使用法: go run main.go [-codeonly] <ディレクトリパス>")
		return
	}
	dirPath := flag.Arg(0)

	var totalLines int64

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			for _, dirToSkip := range skipDirs {
				if info.Name() == dirToSkip {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// --- File processing logic ---

		if *codeOnly {
			// "-codeonly" mode: only count files with specified code extensions
			ext := filepath.Ext(info.Name())
			isCodeFile := false
			for _, codeExt := range codeExtensions {
				if strings.EqualFold(ext, codeExt) {
					isCodeFile = true
					break
				}
			}
			if !isCodeFile {
				return nil // Not a code file, skip
			}
		} else {
			// Default mode: skip specific files and binaries
			for _, fileToSkip := range skipFiles {
				if info.Name() == fileToSkip {
					return nil // Skip this specific file
				}
			}

			isBinary, err := isLikelyBinary(path)
			if err != nil || isBinary {
				return nil // Skip binary files or files that cause errors
			}
		}

		// Count lines for the selected file
		lineCount, err := countLines(path)
		if err != nil {
			fmt.Printf("エラー: %s の行数をカウントできませんでした - %v\n", path, err)
			return nil
		}
		totalLines += lineCount

		return nil
	})

	if err != nil {
		fmt.Printf("ディレクトリのスキャン中にエラーが発生しました: %v\n", err)
		return
	}

	fmt.Printf("合計行数: %d\n", totalLines)
}

func isLikelyBinary(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false, err
	}

	for i := 0; i < n; i++ {
		if buffer[i] == 0 {
			return true, nil
		}
	}
	return false, nil
}

func countLines(filePath string) (int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var count int64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		count++
	}
	return count, scanner.Err()
}
