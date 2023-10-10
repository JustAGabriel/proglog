package internal

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path"
)

func randDirName() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func testRootDir() string {
	return path.Join("/", "tmp")
}

// GetTempFile creates a file in the temporary 'dir'.
// If 'dir' is "", a random dir name will be assigned.
func GetTempFile(dir, filename string) *os.File {
	// make sure that a random dir is generated if none specified
	if dir == "" {
		dirName := randDirName()
		dir = GetTempDir(dirName)
	} else {
		dir = path.Join(testRootDir(), dir)
	}

	f, err := os.CreateTemp(dir, filename)
	if err != nil {
		panic(err)
	}
	return f
}

// GetTempDir creates a temporary directory 'dir'.
// if the 'dir' string includes '*', those will be replaced with a random string.
func GetTempDir(dir string) string {
	dirPath, err := os.MkdirTemp(testRootDir(), dir)
	if err != nil {
		panic(err)
	}
	return dirPath
}
