package internal

import (
	"os"
)

func GetTempFile(dirName, filename string) *os.File {
	f, err := os.CreateTemp(dirName, filename)
	if err != nil {
		panic(err)
	}
	return f
}

func GetTempDir(dirName string) string {
	dirPath, err := os.MkdirTemp("", dirName)
	if err != nil {
		panic(err)
	}
	return dirPath
}
