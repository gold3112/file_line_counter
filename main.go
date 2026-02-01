package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("使用法: go run main.go <ディレクトリパス>")
		return
	}
	dirPath := os.Args[1]

	var totalLines int64

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			isBinary, err := isLikelyBinary(path)
			if err != nil {
				fmt.Printf("警告: %s のファイルタイプを判別できませんでした - %v\n", path, err)
				return nil
			}

			if !isBinary {
				lineCount, err := countLines(path)
				if err != nil {
					fmt.Printf("エラー: %s の行数をカウントできませんでした - %v\n", path, err)
					return nil // 他のファイルの処理を続行
				}
				totalLines += lineCount
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("ディレクトリのスキャン中にエラーが発生しました: %v\n", err)
		return
	}

	fmt.Printf("合計行数: %d\n", totalLines)
}

// isLikelyBinary checks if a file is likely a binary file by searching for null bytes.
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
			return true, nil // Null byte found, likely binary
		}
	}

	return false, nil // No null bytes found, likely text
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
