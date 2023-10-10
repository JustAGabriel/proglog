package internal

import (
	"os"
	"path"
)

func GetTempFile(dirName, filename string) *os.File {
	dir := path.Join("tmp", dirName)
	f, err := os.CreateTemp(dir, filename)
	if err != nil {
		panic(err)
	}
	return f
}

func GetTempDir(dirName string) string {
	dirPath, err := os.MkdirTemp("tmp", dirName)
	if err != nil {
		panic(err)
	}
	return dirPath
}
